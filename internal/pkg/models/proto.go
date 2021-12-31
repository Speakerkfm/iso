package models

type ProtoPlugin struct {
	ModuleName    string
	Imports       []string
	ProtoServices []*ProtoServiceDesc
}

type ProtoFile struct {
	Name         string
	PkgName      string
	OriginalPath string
	Path         string

	RawData []byte
}

type ProtoServiceDesc struct {
	Name      string
	Methods   []*ProtoMethodDesc
	ProtoPath string
	PkgName   string
}

type ProtoMethodDesc struct {
	Name         string
	RequestType  string
	ResponseType string
}
