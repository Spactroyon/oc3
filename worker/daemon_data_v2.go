package worker

import "fmt"

type (
	daemonDataV2 struct {
		data data
	}

	data map[string]any
)

func (d data) getString(key string) (s string, err error) {
	if i, ok := d[key]; !ok {
		err = fmt.Errorf("invalid schema: no key '%s'", key)
	} else if s, ok = i.(string); !ok {
		err = fmt.Errorf("invalid schema: %s not a string", key)
	} else if len(s) == 0 {
		err = fmt.Errorf("unexpected empty %s value", key)
	}
	return
}

func (d data) getDictSubKeys(key string) (l []string, err error) {
	m, ok := d[key].(map[string]any)
	if !ok {
		err = fmt.Errorf("invalid schema: no key '%s'", key)
		return
	}
	for s := range m {
		l = append(l, s)
	}
	return
}

func (d *daemonDataV2) nodeNames() (l []string, err error) {
	return d.data.getDictSubKeys("nodes")
}

func (d *daemonDataV2) objectNames() (l []string, err error) {
	return d.data.getDictSubKeys("services")
}

func (d *daemonDataV2) clusterID() (s string, err error) {
	return d.data.getString("cluster_id")
}

func (d *daemonDataV2) clusterName() (s string, err error) {
	return d.data.getString("cluster_name")
}

func (d *daemonDataV2) parseNodeFrozen(i any) string {
	switch v := i.(type) {
	case int:
		if v > 0 {
			return "T"
		}
	case float64:
		if v > 0 {
			return "T"
		}
	}
	return "F"
}

func (d *daemonDataV2) nodeFrozen(nodename string) (s string, err error) {
	i, ok := mapTo(d.data, "nodes", nodename, "frozen")
	if !ok {
		err = fmt.Errorf("can't retrieve frozen for %s", nodename)
		return
	} else {
		return d.parseNodeFrozen(i), nil
	}
}

func mapTo(m map[string]any, k ...string) (any, bool) {
	if len(k) <= 1 {
		v, ok := m[k[0]]
		return v, ok
	}
	if v, ok := m[k[0]].(map[string]any); !ok {
		return v, ok
	} else {
		return mapTo(v, k[1:]...)
	}
}

func (d *daemonDataV2) getFromKeys(keys ...string) (v any, err error) {
	if v, ok := mapTo(d.data, keys...); !ok {
		return v, fmt.Errorf("getFromKeys can't expand from %v", keys)
	} else {
		return v, nil
	}
}
