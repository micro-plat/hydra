package conf

type Server struct {
	Proto string `json:"proto" valid:"ascii,required"`
}
