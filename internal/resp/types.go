package resp

type Typ struct {
	RespCode byte
	Name string
}

var (
	ERROR   = Typ { RespCode: '-', Name: "Error" }
	STRING  = Typ { RespCode: '+', Name: "Error" }
	INTEGER = Typ { RespCode: ':', Name: "Error" }
	BULK    = Typ { RespCode: '$', Name: "Error" }
	ARRAY   = Typ { RespCode: '*', Name: "Error" }
)
