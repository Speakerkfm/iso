package models

const (
	PluginName = "struct.so"
)

// ProtoPluginDesc сущность, которая хранит объекты для генерации прото плагина
type ProtoPluginDesc struct {
	ModuleName    string
	Imports       []string
	ProtoServices []*ProtoServiceDesc
}

// ProtoFile сущность, которая хранит данные .proto файла и его содержимое
type ProtoFileData struct {
	Name         string
	PkgName      string
	OriginalPath string
	Path         string

	RawData []byte
}

// ProtoServiceDesc сущность, которая содержит описание прото сервиса из .proto файла
type ProtoServiceDesc struct {
	Name      string
	Methods   []*ProtoMethodDesc
	ProtoPath string
	PkgName   string
}

// ProtoMethodDesc сущность, которая содержит описание прото метода из .proto файла
type ProtoMethodDesc struct {
	Name         string
	RequestType  string
	ResponseType string
}
