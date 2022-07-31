package profession

import (
	"github.com/eolinker/eosc"
	"gopkg.in/yaml.v2"
)

func NewProfessions(data []byte) (*Professions, error) {
	p := &Professions{}
	err := p.Reset(data)
	return p, err
}

type Professions struct {
	configs []*eosc.ProfessionConfig
}

func (p *Professions) Reset(data []byte) error {
	var configs map[string][]*eosc.ProfessionConfig
	err := yaml.Unmarshal(data, &configs)
	if err != nil {
		return err
	}
	if v, ok := configs["professions"]; ok {
		p.configs = v
	}
	return nil
}

func (p *Professions) Sort() []*eosc.ProfessionConfig {
	list := p.configs
	sl := make([]*eosc.ProfessionConfig, 0, len(list))
	sm := make(map[string]int)
	index := 0
	for i, p := range list {
		if p.Mod == eosc.ProfessionConfig_Singleton {
			sl = append(sl, p)
			sm[p.Name] = index
			index++
			list[i] = nil
		}
	}
	for len(list) > 0 {
		sc := 0
		for i, v := range list {
			if v == nil {
				sc++
				continue
			}
			dependenciesHas := 0
			for _, dep := range v.Dependencies {
				if _, has := sm[dep]; !has {
					break
				}

				dependenciesHas++
			}
			if dependenciesHas == len(v.Dependencies) {
				sl = append(sl, v)
				sm[v.Name] = index
				index++
				sc++
				list[i] = nil
			}
		}
		if sc == 0 {
			// todo profession dependencies error
			break
		}
		tmp := list[:0]
		for _, v := range list {
			if v != nil {
				tmp = append(tmp, v)
			}
		}
		list = tmp
	}
	return sl
}
