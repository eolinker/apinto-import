package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apinto-import/request"

	//"gopkg.in/yaml.v2"
	"github.com/eolinker/apinto-import/profession"
	"github.com/ghodss/yaml"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/utils/zip"
	"github.com/urfave/cli/v2"
)

func Import() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: fmt.Sprintf("import %s server", env.AppName()),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "apinto-address",
				Aliases:  []string{"addr"},
				Usage:    "like this{scheme}://{ip}:{port}",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "file path",
				Required: true,
			},
		},
		Action: ImportFunc,
	}
}

func ImportFunc(c *cli.Context) error {
	addr := c.String("addr")
	path := c.String("path")
	if path == "" {
		return errors.New("invalid path")
	}
	unzipPath := strings.TrimSuffix(path, ".zip")
	err := zip.DeCompress(path, unzipPath)
	if err != nil {
		return err
	}
	paths, err := getPathsOnDir(unzipPath, "profession-")
	if err != nil {
		return err
	}
	cfgs, err := getData(paths)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/professions.yml", unzipPath))
	if err != nil {
		return err
	}

	p, err := profession.NewProfessions(data)
	if err != nil {
		return err
	}

	ps := p.Sort()
	for _, s := range ps {
		if cfg, ok := cfgs[s.Name]; ok {
			if len(cfg) < 1 {
				continue
			}
			uri := fmt.Sprintf("%s/api/%s", addr, s.Name)
			for _, v := range cfg {
				body, err := json.Marshal(v)
				if err != nil {
					return fmt.Errorf("fail to post %s data,error is %w", s.Name, err)
				}
				_, err = request.PostData(uri, body)
				if err != nil {
					return fmt.Errorf("fail to post %s data,error is %w,body is %s", s.Name, err, string(body))
				}
			}
		}
	}

	log.Info("import successful")
	return nil
}

type Response struct {
	StatusCode int    `json:"status"`
	Data       []byte `json:"data"`
}

func getPathsOnDir(dir string, filePrefix string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("the path is not dir")
	}
	files, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			ps, err := getPathsOnDir(fmt.Sprintf("%s/%s", dir, file.Name()), filePrefix)
			if err != nil {
				return nil, err
			}
			paths = append(paths, ps...)
			continue
		}
		if strings.HasPrefix(file.Name(), filePrefix) {
			paths = append(paths, fmt.Sprintf("%s/%s", dir, file.Name()))
		}
	}
	return paths, nil
}

func getData(paths []string) (map[string][]interface{}, error) {
	cfgs := make(map[string][]interface{})
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var cfg map[string][]interface{}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return nil, err
		}
		for key, value := range cfg {
			if _, ok := cfgs[key]; !ok {
				cfgs[key] = make([]interface{}, 0, len(value))
			}
			cfgs[key] = append(cfgs[key], value...)
		}
	}
	return cfgs, nil
}
