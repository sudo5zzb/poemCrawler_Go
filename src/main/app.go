package main
import(
	"fmt"
	"../config"
)

var conf *config.Config

func init(){
	conf=config.GetConfig()
}

func main()  {
	fmt.Println(conf)
}