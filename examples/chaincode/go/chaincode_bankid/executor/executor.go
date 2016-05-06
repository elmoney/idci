package main

import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
)

func postreq(reqType int) {

	var (
		url string
		jsonStr []byte
		address = "http://localhost:5000"
	)

	switch reqType {

	case 1:
		// login
		url = address + "/registrar"
		jsonStr = []byte(`{
			  "enrollId": "lukas",
  				"enrollSecret": "NPKYL39uKbkj"
			}`)

	case 2:
		// init chaincode
		url = address + "/devops/deploy"
		jsonStr = []byte(`{
				  "type": "GOLANG",
				  "chaincodeID":{
				      "name":"bicc"
				  },
				  "ctorMsg": {
				      "function":"init",
				      "args":["a", "100", "b", "200"]
				  },
				  "secureContext": "lukas"
			}`)

	case 10:
		// add webmoney public key
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
				"chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
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
			}`)

	case 11:
		// add test-bank public key
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
			"chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
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
			}`)

	case 20:
		// get webmoney public key
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-pubkey",
					          "args":["webmoney"]
					      },
					      "secureContext": "lukas"
					  }
		}`)

	case 21:
		// get test-bank public key
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-pubkey",
					          "args":["test-bank"]
					      },
					      "secureContext": "lukas"
					  }
				}`)

	case 30:
		// create request webmoney -> test-bank
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
			  "chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
			      },
			      "ctorMsg":{
			          "function":"create-request",
			          "args":[
			          	"08000000536c125e8ef944e8be34419d5bab4245f8369b609c1ce71f59c2475b0f50b16b4664ebf385dd9b64b7fc8099d6858c88bf09b9b49998f28e858d02fa8620c73ee778bdff180b7231d6047ed970fd27cf3af31caa6b2b16f10f1e2e2b33f7c91af46889d9d2268652f2e02a657b6c3b15b2562d37cd6db7498afa215b86005122ca68b53e1de4403dbd5c336feffc3cb72a9e0ee2e218c132c1c7d0aff4e01798d23a5f9721ed729cfb9bf2d8ab8f320fb475c4b05e2e80cdab17b49e85c1458548c4231f94ba292e1482205e19b1d01e653bfc8841a5851146a73e70e3d0d04b9d4e0a30367c3ab516520f2c95c57f026a717b88ff8f333926cfb83fdb1ec1d8",
			          	"3f31f7851874e064b92694ba6ad1d5441e380cea06105e57f945dccd347f71d569ac0afdf2615519741a8091ef5fcff275233d7660e018aa5cf7a334082f58e869d9545379c49881da7e3becc1bb2ab4faabea36928df3c913a2901359d035ddb29f90aea947a4d610efd5de95440d034a2ed3aaff07dcc3b221af2f22ec5969f5d8caf278c29018f7d1f603f0099f32ee9292a5d0d89c148e2c908040e94e91e9d4097998c09446af4795dbcb6ded8368db58d002ce9254edb18ead8ecb37db7c40d6138b8bd27db5ef17d2478c364f89b1c07d5ffe31373e79523848e3cc23fec4b466af18204a8ab140781b4f8fa20d8b1ad3297a787eafc41c668a194445",
			          	"15909585-16f0-406a-940a-df49cefbd652",
				      	"0900000008b3f2c498b63c3a08f359ec405ed1f615f1ced97cc5a7c1b86d4303305022dd2e29c3baec037c5191486f1e24d0d98a5cec9dfdd8539b2edb32ac14b9b633a4354063563a83d629caa2f24244a413191d16b401a923b790957972586e783feea71eb33e6b13611555f5a024b2405952e4d0a54b3c3e33c8c127cd2044710da93fbffa9c0c119a1d5daf1e4e38b4e03ad49a2d3baea71959e88f1bc12cbe3d4a1beb498d1e299878f79a50058961c9430dc5fdb29053d018257c66e843a0184ec7b0b48bda106b69a95706132ec1f372d4f8f3d72db888417e6f87a5da4154d0e220f9f5f45fc84cb4d0ef0d3f16da7c9963e93bd5941a9db40587718bdbbc4e",
				      	"010000003798406f889ad1d47012ac1adf0c7885374a7e1510c318a8321e517a538ebd692fcefa4d9ae49f68f280a702c651c9dfb8bc57dccdf1084c1e3d42ae9d3aa69796f7b7abd4c5ef7bf4ea484851de646892395754f6376ac451ad18d93fa824b4715158ff4ecfab280757d9d7d972c56c90f34c75f20a101e49e22c49787ff49d472f6d8057e4bbf6bd752e12545a5336911e770a079f768c291816201bf4fdb4585f90790240b5497cc3a0f1a52cfe21923cb6fcc95dbd5dbfee1c45c0cb38a2ac799a5f7126943276a34f8d902956d72e76ad6d7de4d990cc17f4f31a19dc989dba2e9ee7dd93b50c82dbd8e95cf32722f79fefb19a2f4528ecedc5e6c31962",
				      	"0c00000080c02b5ed44f7f748f358655e4978f811b46b68862ce762dcf546e310b1f92b64ebb2d8736390e1d4a66bb92980c42ea9ffe1b9593a3ca47f2174dca69d6bd6a56a0c4f0beaa6e288282fca5d9a19fe4f4f0a91c1137abf8debe075581b285122af75f45450684239961c3aed2fb6d28f2b4ffc5e20cef998f91ece742a64e360afd342f02b3f1fb38d6bf7cb84cca1c40ac9e7110b00ef0dbd5ea9117bcee1bbb605b86635c584fc5f752c55166e1673eaa020a8420094119a8de5dc8733492efb28a94c70922044723c136b9efac9f6b78f7f15e61fad5163658ea1c9f4a8782eeb2eee63becba9c87be61ba937235eb6089f8fdabd88c5b72dc87195c1cb0",
				      	"0100000029c0e1b500c9f93bf8256e641350b0916e1579be8bd8ff832578db495f4a0851971b4d87cdd3a0ec84568a929f09950ee619b366bea0228a3ac625992e9249bb27ade25a76fba350d5a0589a1a4d2567eaeeef1667f79064f62a39183e73b5beac247694f902efc2c7c4cdb9770153d311e4def7f569d8107edff3752cb7ff17e36794ed5473a95a48096e704152aca1d828505daad6b3bbd72b86c971f2082d529a5067d8fc602bacc6669882e4b050f56843ea9c27394a914e8c8196b38aa1e1f47703f251344a745327164748a4037382cc0d63ca90ceebd8a33e71b8fa9f291c9f2542cb21ac66e260be5c88ad7105cf31654c181a10cb97cf4e7f532650",
				      	"340000008380e3d91a831c2591a371d0f2cd33f2045e35e49db1e65ccc1bb364bf052bb5d5776aa579dbb59befc0e4b2b2f3cf7d8b21a2868f1938a8f2c755b68ecf8fbddc19b2163a4bb8993b278a2cf419e59f2b2fbc23fed5dcf724d6d22b75a14255f499143ae1213cabcd3ee5bc1d1dee53225fa3da21c76bde8e2bf56a2eb5a5c2e2a7b97c56b7a30d60eb5323dd5baf1b76fbbe8b64f63057201a01406bcc01e4ec5fa1076bbbc64e0a424f94cc0a2843372c751b6bce0798ba656ebf77fa6a39c7c7dae3848b371d9f5c150c586e06572cfd335c30da346c0275e3bd550bab9d8432550e406c1383b1b9ee5530bf1934f15f5481fdd32fccb2964078f31bc379"
				      ]
			      },
			      "secureContext": "lukas"
			  }
			}`)

	case 31:
		// create request webmoney -> test-bank
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
			  "chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
			      },
			      "ctorMsg":{
			          "function":"create-request",
			          "args":[
			          	"08000000536c125e8ef944e8be34419d5bab4245f8369b609c1ce71f59c2475b0f50b16b4664ebf385dd9b64b7fc8099d6858c88bf09b9b49998f28e858d02fa8620c73ee778bdff180b7231d6047ed970fd27cf3af31caa6b2b16f10f1e2e2b33f7c91af46889d9d2268652f2e02a657b6c3b15b2562d37cd6db7498afa215b86005122ca68b53e1de4403dbd5c336feffc3cb72a9e0ee2e218c132c1c7d0aff4e01798d23a5f9721ed729cfb9bf2d8ab8f320fb475c4b05e2e80cdab17b49e85c1458548c4231f94ba292e1482205e19b1d01e653bfc8841a5851146a73e70e3d0d04b9d4e0a30367c3ab516520f2c95c57f026a717b88ff8f333926cfb83fdb1ec1d8",
			          	"9b557ffa11d44ac957023b8b9a2d9eb6753ef6697f204e576b1b3a63fab06e2a814794f31fe8272c1b9147155a1449e9ca3e54d4882ce07a0eb3e4cfdef3b3eb632f68769dcaad0a45246b36de23c0ec9dfc997696b296c8bb361428410f2fec3983e8702fbebe0125bd7d5b559b9d12bc402b57edcaf43c126cbbed33c1755d1f0efd48eddcf1a9fd1e59f93c322e35c82b30dd6498e92305c97c64350e8cbdb1f5a983ce3e944cde0cfd9c88825746e131a291783103e1f962bee7765b307fcf9e1c664fb56b8602bddf7ebea8ab533b59cfbc46408dd9c83f7dd0238953be0b68e56b14149408fcbe209b6f2c7ba999a61d7b6c7f03e49f5279197acc073d",
			          	"847f05ce-e37c-4036-810e-9bddb7af659f",
				      	"0900000008b3f2c498b63c3a08f359ec405ed1f615f1ced97cc5a7c1b86d4303305022dd2e29c3baec037c5191486f1e24d0d98a5cec9dfdd8539b2edb32ac14b9b633a4354063563a83d629caa2f24244a413191d16b401a923b790957972586e783feea71eb33e6b13611555f5a024b2405952e4d0a54b3c3e33c8c127cd2044710da93fbffa9c0c119a1d5daf1e4e38b4e03ad49a2d3baea71959e88f1bc12cbe3d4a1beb498d1e299878f79a50058961c9430dc5fdb29053d018257c66e843a0184ec7b0b48bda106b69a95706132ec1f372d4f8f3d72db888417e6f87a5da4154d0e220f9f5f45fc84cb4d0ef0d3f16da7c9963e93bd5941a9db40587718bdbbc4e",
				      	"010000003798406f889ad1d47012ac1adf0c7885374a7e1510c318a8321e517a538ebd692fcefa4d9ae49f68f280a702c651c9dfb8bc57dccdf1084c1e3d42ae9d3aa69796f7b7abd4c5ef7bf4ea484851de646892395754f6376ac451ad18d93fa824b4715158ff4ecfab280757d9d7d972c56c90f34c75f20a101e49e22c49787ff49d472f6d8057e4bbf6bd752e12545a5336911e770a079f768c291816201bf4fdb4585f90790240b5497cc3a0f1a52cfe21923cb6fcc95dbd5dbfee1c45c0cb38a2ac799a5f7126943276a34f8d902956d72e76ad6d7de4d990cc17f4f31a19dc989dba2e9ee7dd93b50c82dbd8e95cf32722f79fefb19a2f4528ecedc5e6c31962",
				      	"0c00000080c02b5ed44f7f748f358655e4978f811b46b68862ce762dcf546e310b1f92b64ebb2d8736390e1d4a66bb92980c42ea9ffe1b9593a3ca47f2174dca69d6bd6a56a0c4f0beaa6e288282fca5d9a19fe4f4f0a91c1137abf8debe075581b285122af75f45450684239961c3aed2fb6d28f2b4ffc5e20cef998f91ece742a64e360afd342f02b3f1fb38d6bf7cb84cca1c40ac9e7110b00ef0dbd5ea9117bcee1bbb605b86635c584fc5f752c55166e1673eaa020a8420094119a8de5dc8733492efb28a94c70922044723c136b9efac9f6b78f7f15e61fad5163658ea1c9f4a8782eeb2eee63becba9c87be61ba937235eb6089f8fdabd88c5b72dc87195c1cb0",
				      	"0100000029c0e1b500c9f93bf8256e641350b0916e1579be8bd8ff832578db495f4a0851971b4d87cdd3a0ec84568a929f09950ee619b366bea0228a3ac625992e9249bb27ade25a76fba350d5a0589a1a4d2567eaeeef1667f79064f62a39183e73b5beac247694f902efc2c7c4cdb9770153d311e4def7f569d8107edff3752cb7ff17e36794ed5473a95a48096e704152aca1d828505daad6b3bbd72b86c971f2082d529a5067d8fc602bacc6669882e4b050f56843ea9c27394a914e8c8196b38aa1e1f47703f251344a745327164748a4037382cc0d63ca90ceebd8a33e71b8fa9f291c9f2542cb21ac66e260be5c88ad7105cf31654c181a10cb97cf4e7f532650",
				      	"340000008380e3d91a831c2591a371d0f2cd33f2045e35e49db1e65ccc1bb364bf052bb5d5776aa579dbb59befc0e4b2b2f3cf7d8b21a2868f1938a8f2c755b68ecf8fbddc19b2163a4bb8993b278a2cf419e59f2b2fbc23fed5dcf724d6d22b75a14255f499143ae1213cabcd3ee5bc1d1dee53225fa3da21c76bde8e2bf56a2eb5a5c2e2a7b97c56b7a30d60eb5323dd5baf1b76fbbe8b64f63057201a01406bcc01e4ec5fa1076bbbc64e0a424f94cc0a2843372c751b6bce0798ba656ebf77fa6a39c7c7dae3848b371d9f5c150c586e06572cfd335c30da346c0275e3bd550bab9d8432550e406c1383b1b9ee5530bf1934f15f5481fdd32fccb2964078f31bc379"
				      ]
			      },
			      "secureContext": "lukas"
			  }
			}`)

	case 40:
		// approve cr request
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
			  "chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
			      },
			      "ctorMsg":{
			          "function":"approve-request",
			          "args":[
				      	"0900000041ea8a031944098c5ac0ba3c73065bb2f5b44424045ad8d46db4b42bb88c478c842e6d32efde82b2a7daa712049f1fe1e31bb8bd5cd0eccb6ddd60ff2eec5fa0bd20dbfb3ae50d1ef71202bf2311242cbbde03c3cf4ecee2a5a1b29410d35aa99cf0909724618aea0fe909259165e4245c1fd87a6a5339e00268c3e6cdc3fcdfd3770dda1703b262803aee393c4cba5e9cbff7401929503ba3a5e59aceb5769602291318f971c76ded0fabb1cb4a6684c5456ea23b3daec6fc97585d7c20f8bcf11dfa381084d17d91445354b7559059087d737adc48e6a99f30ddcef2e7b75a7e79033cacea0f7fc7987f2975457f65f8f37567133f9761cf89ee304fbe1190",
				      	"875f505520609b9cf5eaac85fdfe2a48edb9d01cabdc4c4ab4efbabcca4438dfcbf608a716dda12240bb350e520dfa8071e87debedcdce5b40163a1fbb93c3dc47936b46d7433426d633fc74adbc0aef391ec7b8bb25ea10f0348040c2c7608c71147e7dd3bdd0a716ba7124f37f8447f861721b95a5621c85b8cc6d2418627d2b85ab7027060a011ee553fee3e5e964810fd70b8fd50d42aa1994764689302cac715848c457682e65ae235d1d9136320342ae5ae459675ea6a7b0c39c10f69b6a47757ddc4b0df8f6afc27d0b6c93c2eb4628ff62a759746ad0da7fa8dabe7f52f0412cc53229d24c59f4a80095e9d3e3d748fb39f2d34221090623c20a6f40",
				      	"15909585-16f0-406a-940a-df49cefbd652",
				      	"070000000673b4f148f73a74ecc52f4c47851c85c6f9a39cc38c9e85b9609162592f8d1471b82ad679581e8a5a4476e471cf3f7a02a260f345a1640f00cd2433a3f4104a25d54b500d33d6e4bad707deee814f02a8350140cf9b170c743f0544788c6d3b1e8184c828b025eeaa1dbf6e4a5e5fabc69c66a5c16a08dbe565f3ebe4168b7b9621030a92e6754f240b03bd2f1e206ea9f221aacff815938c37381a063c81e076f57e6a74603e0ef488e67a3fa0d2ed2aa2a0241bdddadac61c32502dcd1be985d1cdea1aca935d3ffcd45e530a84dc049381728c5a1ebb8d623ee243d277670186d1463cac60c6b7b9935a3c9b69b20ceb17e0ea4358b12603954666bc9d4f"
				      ]
			      },
			      "secureContext": "lukas"
			  }
			}`)

	case 41:
		// reject cr request
		url = address + "/devops/invoke"
		jsonStr = []byte(`{
			  "chaincodeSpec":{
			      "type": "GOLANG",
			      "chaincodeID":{
			          "name":"bicc"
			      },
			      "ctorMsg":{
			          "function":"reject-request",
			          "args":[
				      	"0900000041ea8a031944098c5ac0ba3c73065bb2f5b44424045ad8d46db4b42bb88c478c842e6d32efde82b2a7daa712049f1fe1e31bb8bd5cd0eccb6ddd60ff2eec5fa0bd20dbfb3ae50d1ef71202bf2311242cbbde03c3cf4ecee2a5a1b29410d35aa99cf0909724618aea0fe909259165e4245c1fd87a6a5339e00268c3e6cdc3fcdfd3770dda1703b262803aee393c4cba5e9cbff7401929503ba3a5e59aceb5769602291318f971c76ded0fabb1cb4a6684c5456ea23b3daec6fc97585d7c20f8bcf11dfa381084d17d91445354b7559059087d737adc48e6a99f30ddcef2e7b75a7e79033cacea0f7fc7987f2975457f65f8f37567133f9761cf89ee304fbe1190",
				      	"956027d4092f4a7f8f8422af9548d7531dedebc158079c6ba948f567822ad84007cd6223f1f8883bd8839210dca8bf2f7b4592762be55250f9cfd2542f123bc3f19becf24a2e8aa07ab408b4bf4d1acecdf4261967f462a633dfb605e8a697fb91121638e936ab918e5398fbbb943b88101a2d11c0efc62714bd59c8bf123694a70dee17051aee267e10b15a864adc0c6f93a2a61fed8eee09fd4ea5b7ee97b2f65854a8be9829cf6c987f79427cb52fbcd373953c371477843680bdc752a645a1da26634f9342de634bbb24796707a8903b2522b268e7f422c90ed23d34f8329ab4f9155eb1bab2ca60eab05706de651518523b57d8db7464e720058d89f8a2",
				      	"847f05ce-e37c-4036-810e-9bddb7af659f",
				      	"070000000673b4f148f73a74ecc52f4c47851c85c6f9a39cc38c9e85b9609162592f8d1471b82ad679581e8a5a4476e471cf3f7a02a260f345a1640f00cd2433a3f4104a25d54b500d33d6e4bad707deee814f02a8350140cf9b170c743f0544788c6d3b1e8184c828b025eeaa1dbf6e4a5e5fabc69c66a5c16a08dbe565f3ebe4168b7b9621030a92e6754f240b03bd2f1e206ea9f221aacff815938c37381a063c81e076f57e6a74603e0ef488e67a3fa0d2ed2aa2a0241bdddadac61c32502dcd1be985d1cdea1aca935d3ffcd45e530a84dc049381728c5a1ebb8d623ee243d277670186d1463cac60c6b7b9935a3c9b69b20ceb17e0ea4358b12603954666bc9d4f"
				      ]
			      },
			      "secureContext": "lukas"
			  }
			}`)

	case 50:
		// get request by id
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-request",
					          "args":[
					          	"090000001d9cd6ef4554ddf7780452bd58aa38aeba9f3ac91d28066936d0c3bdc4b3b9ea8a6e42d696d2126d702d12715e46284a4cf80a8a1605a3075be98b3b55041598ea427ba8d19dc0f3f21ac704aef9ab3b2a70ec19ae58b51d03a2c26cb0842ba13419e96fb459054e2f48190a20e878e0f64593b0b514c21646ca170a9d9256bbbc54a961a24e75811135aa6624a79b747ee224f0a6832c7a2a5d892594b898d97bf2d79f6e3ea824e2c83d130ca5ee6fbabca69ca2f0508e5124c5a8761a07fcc13edd0b43da4d90b25eacf57107b16da6ad1f5c9c1df379e51db10eed1c80eb65e4814bfb77a8018ad0b0ba0e4776d54957fa3436953681c5f2585c53c8b524",
					          	"15909585-16f0-406a-940a-df49cefbd652"
					          ]
					      },
					      "secureContext": "lukas"
					  }
				}`)

	case 51:
		// get request by id
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-request",
					          "args":[
					          	"090000001d9cd6ef4554ddf7780452bd58aa38aeba9f3ac91d28066936d0c3bdc4b3b9ea8a6e42d696d2126d702d12715e46284a4cf80a8a1605a3075be98b3b55041598ea427ba8d19dc0f3f21ac704aef9ab3b2a70ec19ae58b51d03a2c26cb0842ba13419e96fb459054e2f48190a20e878e0f64593b0b514c21646ca170a9d9256bbbc54a961a24e75811135aa6624a79b747ee224f0a6832c7a2a5d892594b898d97bf2d79f6e3ea824e2c83d130ca5ee6fbabca69ca2f0508e5124c5a8761a07fcc13edd0b43da4d90b25eacf57107b16da6ad1f5c9c1df379e51db10eed1c80eb65e4814bfb77a8018ad0b0ba0e4776d54957fa3436953681c5f2585c53c8b524",
					          	"847f05ce-e37c-4036-810e-9bddb7af659f"
					          ]
					      },
					      "secureContext": "lukas"
					  }
				}`)

	case 52:
		// get request for approve by unit test-bank
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-requests-for-approve",
					          "args":[
					          	"090000001d9cd6ef4554ddf7780452bd58aa38aeba9f3ac91d28066936d0c3bdc4b3b9ea8a6e42d696d2126d702d12715e46284a4cf80a8a1605a3075be98b3b55041598ea427ba8d19dc0f3f21ac704aef9ab3b2a70ec19ae58b51d03a2c26cb0842ba13419e96fb459054e2f48190a20e878e0f64593b0b514c21646ca170a9d9256bbbc54a961a24e75811135aa6624a79b747ee224f0a6832c7a2a5d892594b898d97bf2d79f6e3ea824e2c83d130ca5ee6fbabca69ca2f0508e5124c5a8761a07fcc13edd0b43da4d90b25eacf57107b16da6ad1f5c9c1df379e51db10eed1c80eb65e4814bfb77a8018ad0b0ba0e4776d54957fa3436953681c5f2585c53c8b524"
					          ]
					      },
					      "secureContext": "lukas"
					  }
				}`)

	case 53:
		// get request for approve by unit test-bank
		url = address + "/devops/query"
		jsonStr = []byte(`{
					  "chaincodeSpec":{
					      "type": "GOLANG",
					      "chaincodeID":{
					          "name":"bicc"
					      },
					      "ctorMsg":{
					          "function":"get-requests-all",
					          "args":[
					          	"090000001d9cd6ef4554ddf7780452bd58aa38aeba9f3ac91d28066936d0c3bdc4b3b9ea8a6e42d696d2126d702d12715e46284a4cf80a8a1605a3075be98b3b55041598ea427ba8d19dc0f3f21ac704aef9ab3b2a70ec19ae58b51d03a2c26cb0842ba13419e96fb459054e2f48190a20e878e0f64593b0b514c21646ca170a9d9256bbbc54a961a24e75811135aa6624a79b747ee224f0a6832c7a2a5d892594b898d97bf2d79f6e3ea824e2c83d130ca5ee6fbabca69ca2f0508e5124c5a8761a07fcc13edd0b43da4d90b25eacf57107b16da6ad1f5c9c1df379e51db10eed1c80eb65e4814bfb77a8018ad0b0ba0e4776d54957fa3436953681c5f2585c53c8b524"
					          ]
					      },
					      "secureContext": "lukas"
					  }
				}`)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error at request doing:\t", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {

	var input string

	for {

		fmt.Println("\t-type 'q'<Enter> to quit")
		fmt.Print("\t-type name of chaincode command<Enter> - ")
		fmt.Scanln(&input)
		fmt.Println("\n")

		if (input == "q") {
			break
		}

		switch input {

		// login
		case "l":
			postreq(1)

		// init cc
		case "ii":
			postreq(2)

		// add unit public key
		case "aw":
			postreq(10)
		case "atb":
			postreq(11)

		// full init
		case "fi":
			postreq(2)
			postreq(10)
			postreq(11)

		// query unit public key
		case "qw":
			postreq(20)
		case "qtb":
			postreq(21)

		// create request 15909585-16f0-406a-940a-df49cefbd652
		case "cr":
			postreq(30)
			//create request 847f05ce-e37c-4036-810e-9bddb7af659f
		case "cr2":
			postreq(31)

		// approve request 15909585-16f0-406a-940a-df49cefbd652
		case "ar":
			postreq(40)
		// reject request 847f05ce-e37c-4036-810e-9bddb7af659f
		case "rr":
			postreq(41)

		// query approve\reject requests
		case "qr":
			postreq(50)
		case "qr2":
			postreq(51)
		case "qrtb":
			postreq(52)
		case "qrtball":
			postreq(53)

		default:
			fmt.Println("Incorrect argument", input)
		}

	}
}
