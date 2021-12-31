package models

type ProtoPlugin struct {
	ModuleName    string
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
}

type ProtoMethodDesc struct {
	Name         string
	RequestType  string
	ResponseType string
}
