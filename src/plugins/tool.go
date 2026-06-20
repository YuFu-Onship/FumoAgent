package plugins

type LangZone map[string]string
type LangPack map[string]LangZone

// {"type":"func","name":"get_weather","args":{"city":"北京"}}

type PluginMeta struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}
