package utils


import (
	"fmt"
	"strings"
	"log"
	"net/http"
)

func QueryStateGet(chainCodeName string,peerQueryUrl string , args ...string) (int,string) {

	postData:=`{
  "chaincodeSpec":{
      "type": "GOLANG",
      "chaincodeID":{
          "name":"@cc"
      },
      "ctorMsg":{
          "function":"query",
          "args":[@key]
      }
  }
} `
	argString:=""
	for _,val := range args {
		argString=argString+fmt.Sprintf(`,"%s"`,val)
	}
	argString=strings.Trim(argString,",")
	postData = strings.Replace(postData,"@key",argString,-1)
	postData = strings.Replace(postData,"@cc",chainCodeName,-1)

	status,resp,err := PostJson(peerQueryUrl, postData)
	if(err!=nil) {
		log.Println("qStateGet.Error : "+err.Error())
		status=http.StatusInternalServerError
	}

	return status,resp
}


func QueryStateWithFunctionGet(chainCodeName string,peerQueryUrl string ,function string , args ...string) (int,string) {

	postData:=`{
  "chaincodeSpec":{
      "type": "GOLANG",
      "chaincodeID":{
          "name":"@cc"
      },
      "ctorMsg":{
          "function":"@function",
          "args":[@key]
      }
  }
} `
	argString:=""
	for _,val := range args {
		argString=argString+fmt.Sprintf(`,"%s"`,val)
	}
	argString=strings.Trim(argString,",")
	postData = strings.Replace(postData,"@key",argString,-1)
	postData = strings.Replace(postData,"@cc",chainCodeName,-1)
	postData = strings.Replace(postData,"@function",function,-1)

	log.Println(postData)

	status,resp,err := PostJson(peerQueryUrl, postData)
	if(err!=nil) {
		log.Println("qStateGet.Error : "+err.Error())
		status=http.StatusInternalServerError
	}

	return status,resp
}


func InvokeCC(function string,peerInvokeUrl string,chainCodeName string,args ...string) (int,string) {

	jsonTmpl :=`{
  "chaincodeSpec":{
      "type": "GOLANG",
      "chaincodeID":{
          "name":"@cc"
      },
      "ctorMsg":{
          "function":"@function",
           "args":[@key]
      }
  }
} `

	argString:=""
	for _,val := range args {
		argString=argString+fmt.Sprintf(`,"%s"`,val)
	}
	argString=strings.Trim(argString,",")
	jsonTmpl = strings.Replace(jsonTmpl,"@key",argString,-1)
	jsonTmpl = strings.Replace(jsonTmpl,"@function",function,-1)
	jsonTmpl = strings.Replace(jsonTmpl,"@cc",chainCodeName,-1)

	log.Println(jsonTmpl)

	status,resp,err := PostJson(peerInvokeUrl,jsonTmpl)

	if(err!=nil) {
		log.Println("InvokeCC.Error : "+err.Error())
		status=http.StatusInternalServerError
	}

	return status,resp
}


func InitCC(function string,peerDeployUrl string,chainCodeName string,args ...string)(int,string) {

	jsonTmpl :=`{  "type": "GOLANG",
				  "chaincodeID":{
				      "name":"@cc"
				  },
				  "ctorMsg": {
				       "function":"@function",
                                       "args":[@key]
				  }
			} `


	argString:=""
	for _,val := range args {
		argString=argString+fmt.Sprintf(`,"%s"`,val)
	}
	argString=strings.Trim(argString,",")
	jsonTmpl = strings.Replace(jsonTmpl,"@key",argString,-1)
	jsonTmpl = strings.Replace(jsonTmpl,"@function",function,-1)
	jsonTmpl = strings.Replace(jsonTmpl,"@cc",chainCodeName,-1)

	log.Println(jsonTmpl)

	status,resp,err := PostJson(peerDeployUrl,jsonTmpl)

	if(err!=nil) {
		log.Println("InvokeCC.Error : "+err.Error())
		status=http.StatusInternalServerError
	}

	return status,resp

}

func InitTestUnits(peerInvokeUrl string,chainCodeName string){

	jsonTmpl :=`{"chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"@cc"
			      },
			      "ctorMsg":{
			          "function":"unit-add",
			          "args":[
				      	"webmoney",
				      	"2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b43415145417a63694a59506f39595a4a41304a5077776866530a5376417675737063724e4f2f6f6d346b4d6a4338664c4d536d47645235384e3637473734494a775870596c53534136686d2f61443637555a52417730364636420a6b4746434c4d71353851316e5531415a7767345672412f314f4e437a5465753252536c6545533032654f486b594d4131537154654d4d4f734b504d77445033630a547535626b594157306d353644514a52576467654a7449453141657a4f773631704a503832387633456d6a484b6c346a686376517145464d2b7050727350342b0a44573068516237364b356a514a4266663672782b59766e7470663144644269676a456b6e6276746c7a3235787046626575734b38537476654b2b7758724473320a38725a794679612b354f594964686154537263627932396c3552736979702f594f4e546369445157306650566273452f477842626a77527735477833626d34770a57514944415141420a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d",
				      	"74237dd1f2c5dcb28cce6073c1768a9c91a0aaf3061438e65c525661a8154b6a3132a3f1bd9154923f00edd00e3f93378da1d00bf72d8cff40747b61bd924e7a932b7310c161a0582c20adb04b558634f94e4debd2b74c1310d89a46dd9c62912e1300db2911e845221fa585de0aa350a9032df60b14f36432a6e1435e91a03afbb15fd17706fbb5b8e474b514f64a5ed1abe3db9c1952226f665914cd84726997103afb8991084fb8dfa110cbb92d796bb2a619ac34aa620e737107d1ec5cdc4fd38722fe327533c6ccc4f30c23d5a61e20deded102b1381feb9e24b3d6eba712905a73829119d6333b8ad342d22fb9200536966603c9f9e2143863d1b0a186"
				      ]
			      },
			      "secureContext": "lukas"
			  	}
			} `

	jsonTmpl = strings.Replace(jsonTmpl,"@cc",chainCodeName,-1)

	log.Println(jsonTmpl)
	_,_,err :=PostJson(peerInvokeUrl,jsonTmpl)


	if(err!=nil) {
		log.Println("InvokeCC.Error : "+err.Error())
		return
	}

	jsonTmpl =`{
			"chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"@cc"
			      },
			      "ctorMsg":{
			          "function":"unit-add",
			          "args":[
				      	"test-bank",
				      	"2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b43415145413553466c3067464263444e33677a743938446e490a36323633417947684436566b4c6c50695272444a6b327331713655374a415254672f73373974564a755539724b63357831767a6737716e376c51785a427971530a53796330417258706d77616162466d4142674c44672f706c584352705a714479617367514637366971524d664f74444b4f734d646171523271766a59307a32570a434b55695067734c4f33376e2f524e4b54565a5865666961794769645761617871446137574a6b347950335033626468504276705955764d7a343433617879490a586a53414b4639763852696c4456426c333539367a2f366f6c70724e512f4a6830703933457a6c3139695a6f74334c4e70735265462f797271575671724976310a706c4134735754712f7a5635646968436d6331695078454e6a706b45617742666f6c7336457952676a3241443761556851334d4e50775a6d546e5444386933330a4f514944415141420a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d",
				      	"125af481d6f77eb09613af5dab5188cabc55e0c8bff449248a86da4e2316d49119962f3498f9c5975ba0032d1ff569d7c155ece0c438a46fa9862ed57c9240ecf7c4262b53b3abadd1c129ed2a6acd1812882b52619574e299e44b02242afe829e6600965db35c6c013e23bd2da7081ac51e25001cd9d6785ec5f82cbf682d872f2c9c38d06a6826032b38c4fea9d77532cfd73e58ebc3ff7376d79f555b379e4958b8fa81b97669972fec4d8ad5f4c5a471126da9ca52a389083e7b5b4f7798234657306db590e44254240ef3d59f529c4dc01499ae552ba9f4655a9a016c92d17f4aafa691f86e0b6a141faccac23758333731ef5df241ffac0072d176c5f3"
				      ]
			      }
			}
			} `
	jsonTmpl = strings.Replace(jsonTmpl,"@cc",chainCodeName,-1)
	log.Println(jsonTmpl)
	_,_,err =PostJson(peerInvokeUrl,jsonTmpl)

	if(err!=nil) {
		log.Println("InvokeCC.Error : "+err.Error())
		return
	}
}
