package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	//"github.com/goinggo/mapstructure"
	"reflect"
	"regexp"
	"strconv"

	//"github.com/gofrs/uuid"
	"strings"
	//"github.com/gofrs/uuid"
)

type SimpleContract struct {
	contractapi.Contract
}

//string转map
func buildStringTomap(cond string) (map[string]interface{}, error) {
	var companyMap map[string]interface{}
	err := json.Unmarshal([]byte(cond), &companyMap)
	if err != nil {
		return nil, err
	}

	return companyMap, nil
}

//Struct转Map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
func buidMapListFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]map[string]interface{}, error) {
	entityList := initlist()

	for resultsIterator.HasNext() {
		var mapResult map[string]interface{}
		queryResponse, err := resultsIterator.Next() //获取迭代器中的每一个值
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(queryResponse.Value, &mapResult)
		if err != nil {
			return nil, err
		}
		delete(mapResult, "docType")

		entityList = listappend(entityList, mapResult)
	}

	return entityList, nil
}

func initlist() []map[string]interface{} {
	return make([]map[string]interface{}, 0, 1)

}
func listappend(entityList []map[string]interface{}, entityMap map[string]interface{}) []map[string]interface{} {
	if len(entityList) == cap(entityList) {
		newList := make([]map[string]interface{}, len(entityList), 2*cap(entityList))
		copy(newList, entityList)
		entityList = newList
	}
	entityList = append(entityList, entityMap)

	return entityList
}

func buildErrorResult(err error) map[string]interface{} {
	result := make(map[string]interface{})
	result["code"] = 500
	result["message"] = err.Error()
	result["stackTrace"] = err
	return result
}

func buildSelector(cond string, docType string) (string, error) {
	var conMap map[string]interface{}

	err := json.Unmarshal([]byte(cond), &conMap)
	if err != nil {
		return "", err
	}

	if len(conMap) == 0 {
		return "", errors.New("选择器参数不能为空")
	}
	_, ConsignorNameNotNull := conMap["ConsignorName"]
	if ConsignorNameNotNull {
		ConsignorName := conMap["ConsignorName"]
		delete(conMap, "ConsignorName")
		conMap["ConsignorInfo.ConsignorName"] = ConsignorName

	}
	_, ConsignorIDNotNull := conMap["ConsignorID"]
	if ConsignorIDNotNull {
		ConsignorID := conMap["ConsignorID"]
		delete(conMap, "ConsignorID")
		conMap["ConsignorInfo.ConsignorID"] = ConsignorID

	}
	_, CarrierNameNotNull := conMap["CarrierName"]
	if CarrierNameNotNull {
		CarrierName := conMap["CarrierName"]
		delete(conMap, "CarrierName")
		conMap["CarrierInfo.CarrierName"] = CarrierName

	}
	_, CarrierIDNotNull := conMap["CarrierID"]
	if CarrierIDNotNull {
		CarrierID := conMap["CarrierID"]
		delete(conMap, "CarrierID")
		conMap["CarrierInfo.CarrierID"] = CarrierID

	}

	_, PayerNameNotNull := conMap["PayerName"]
	if PayerNameNotNull {
		PayerName := conMap["PayerName"]
		delete(conMap, "PayerName")
		conMap["Payer.PayerName"] = PayerName

	}

	_, PayerIDNotNull := conMap["PayerID"]
	if PayerIDNotNull {
		PayerID := conMap["PayerID"]
		delete(conMap, "PayerID")
		conMap["Payer.PayerID"] = PayerID
	}

	_, PayeeNameNotNull := conMap["PayeeName"]
	if PayeeNameNotNull {
		PayeeName := conMap["PayeeName"]
		delete(conMap, "PayeeName")
		conMap["Payee.PayeeName"] = PayeeName

	}

	_, PayeeIDNotNull := conMap["PayeeID"]
	if PayeeIDNotNull {
		PayeeID := conMap["PayeeID"]
		delete(conMap, "PayeeID")
		conMap["Payee.PayeeID"] = PayeeID

	}

	conMap["docType"] = docType

	sele := make(map[string]interface{})
	sele["selector"] = conMap

	mjson, err := json.Marshal(sele)
	if err != nil {
		return "", err
	}

	return string(mjson), nil
}

func isAllNumberWithLength(parmar string, len int) bool {
	matched, _ := regexp.MatchString(`^\d{`+strconv.Itoa(len)+`}$`, parmar)
	return matched
}
func isIdentity(parmar string) bool {
	matched, _ := regexp.MatchString(`^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`, parmar)
	return matched
}

func isPassport(parmar string) bool {
	matched, _ := regexp.MatchString(`^1[45][0-9]{7}$|(^[P|p|S|s]\d{7}$)|(^[S|s|G|g|E|e]\d{8}$)|(^[Gg|Tt|Ss|Ll|Qq|Dd|Aa|Ff]\d{8}$)|(^[H|h|M|m]\d{8,10}$)`, parmar)
	return matched
}

const regIdentity string = `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
const regPassport string = `^1[45][0-9]{7}$|(^[P|p|S|s]\d{7}$)|(^[S|s|G|g|E|e]\d{8}$)|(^[Gg|Tt|Ss|Ll|Qq|Dd|Aa|Ff]\d{8}$)|(^[H|h|M|m]\d{8,10}$)`
const regChinese string = `^[^\x00-\xff]+$`
const regEnglish string = `^[a-zA-Z]+$`
const regLevel string = `^[0-5]$|^999$`
const regBizTypeCode string = `^(1002996|1003997|1003998|1002998|1003999)$`
const regWaybillTypeCode string = `^(1|2)$`
const regVehicleTypeCode string = `^(01|02)$`
const regPlateColorCode string = `^([1-5]|9[1-4]|9)$`
const regTrunkTypeCode string = `^(H[1-3][1-9]|H4[1-7]|H5[1-5]|Q[1-3][1-2]|Z[1-5]1|Z71|D1[1-2]|M1[1-5]|M2[1-2]|N11|T11|T2[1-3]|J1[1-3]|G[1-2][1-9]|G3[1-8]|B[1-3][1-9]|B[1-2]A|B[1-2]B|X99)$`
const regShipZoneCode string = `^(B0[1-9]|B10)$`
const regEnergyType string = `^([A-F]|[L-N]|O|[Y-Z])$`
const regGoodsTypeCode string = `^([1-9]00|1[0-5]00|160[1-2]|170[0-1])$`
const regSettlementMeansCodeAndObjectTypeCode string = `^([1-4])$`
const regPaymentMeansCode string = `^([0-4]|[3-4][1-2]|9)$`
const regOrgTypeCode string = `^[1-8]$`
const regAreaCode string = `^[1-8][0-7]\d{4}$`
const regCompanyCode string = `^[A-Z0-9]{18}$`

func regIdentityOrPassword() string {
	return "(" + regIdentity + ")|(" + regPassport + ")"
}
func regAllNumberWithLength(len int) string {
	return `^\d{` + strconv.Itoa(len) + `}$`
}
func regChineseOrEnglish() string {
	return "(" + regChinese + ")|(" + regEnglish + ")"
}
func regWebsite() string {
	return `(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`
}
func regTelPhone() string {
	return `^((13[0-9])|(14[0-9])|(15[0-9])|(17[0-9])|(18[0-9]))\d{8}$`
}
func regDriverLicense() string {
	return `^[A-C][1-2]`
}

const regPlate string = `^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[A-Z]{1}[A-Z0-9]{4}[A-Z0-9挂学警港澳]{1}$`
const regShip string = `^CN[0-9]{11}$`

func regPlateShip() string {
	return "(" + regPlate + ")|(" + regShip + ")"
}

var BankCodeList = []string{"BKCH", "ABOC", "ICBK", "PCBC", "PSBC", "COMM", "CMBC", "SZDB", "CIBK", "MSBC", "GDBK", "SPDB", "FJIB", "EVER", "ZJCB", "CTBA", "ANBK", "ANZB",
	"NOLA", "CHBH", "BRCB", "BJCN", "BARC", "BARB", "KRED", "BTCB", "CDRC", "CBOC", "CSCB", "DDBK", "CRLY", "DGCB", "UOVB", "DLCB", "DSBA", "DECC", "BEAS", "DYSH", "DEUT", "FSBC",
	"KCCB", "BNPA", "NATX", "SGCL", "PNBP", "FZCB", "FXBK", "GDHB", "RCNH", "RCCS", "SDBC", "GLBK", "CZNB", "BGBK", "GYCB", "GZRC", "BKGZ", "GZCB", "BKSH", "BKHD", "HCCB", "HSBC",
	"KODB", "SHBK", "WHCB", "INGB", "RABO", "DRCB", "EWBK", "CITI", "OCBC", "HASE", "CNMB", "HXBK", "HNBN", "HRCN", "HZCB", "JHCB", "BOJJ", "HCRR", "JTCB", "JLBK", "NOSC", "ROYC",
	"BKJN", "DYCB", "HRCB", "HMCB", "RCWJ", "JYCB", "JRCB", "NRCB", "PZCB", "RRCB", "RGCB", "JSCB", "WJRB", "JYHR", "BOJS", "JSHB", "URCC", "YZRC", "JRZC", "BOJX", "CIYU", "JZCB",
	"CKLB", "BOLF", "LJBC", "LSCC", "LWCB", "LYCB", "LCOM", "BOLY", "LZCB", "LZBK", "CHAS", "BOFO", "MSBK", "NAGO", "BOFA", "MBBE", "BOFM", "DHBC", "BINH", "BKNB", "NCCC", "NCCK",
	"NJCB", "HSSY", "YCCB", "NYCB", "BNYM", "SSVB", "BKKB", "QDRC", "QCCB", "BOXN", "JNSH", "ZBBK", "IBKO", "QRCB", "BKQZ", "ESSE", "CRES", "UBSW", "MHCB", "RZCB", "SGRB", "MBTC",
	"ZDRB", "ABNA", "SHRC", "BOSH", "SYCB", "SMBC", "BOTK", "IBXH", "RCXM", "CBXM", "YDBK", "SXCB", "SRCC", "BHRC", "TJRC", "TCCB", "BOTS", "TZBK", "WRCB", "WFCB", "WRCC", "KOEX",
	"WHBK", "UCCB", "HFCB", "WZCB", "IXAB", "CBNB", "DBSS", "COXI", "WPAC", "YBRC", "BCIT", "BROM", "WIHB", "YKCB", "HVBK", "ZFRB", "CRCK", "SCBL", "ICBC", "EIBC", "ADBN", "CTCB",
	"ZRCB", "CZCB", "HCRC", "YHBK", "ZJMT", "ZJTL", "LCBK", "ZXBK", "YWBK", "ZJRB", "ZJPT", "CQRB", "CTGB", "CQCB", "CHCC", "ZXBC", "ZZBK", "9999"}

var InsuranceCompanyCodeList = []string{"ABIC", "AICS", "BOCI", "BPIC", "CAIC", "CCIC", "CICP", "CPIC", "DBIC", "GPIC", "HAIC", "HTIC", "MACN", "PAIC", "PICC", "TAIC", "TPIC",
	"YDCX", "YGBX", "ZKIC", "MACN", "YAIC", "TPBX", "ACIC", "DHIC", "ALIC", "QITA"}

var AreaCodeList = []string{"4", "248", "8", "12", "16", "20", "24", "660", "10", "28", "32", "51", "533", "36", "40", "31", "44", "48", "50", "52", "112", "56", "84", "204",
	"60", "64", "68", "535", "70", "72", "74", "76", "86", "96", "100", "854", "108", "132", "116", "120", "124", "136", "140", "148", "152", "162", "166", "170", "174", "180",
	"178", "184", "188", "384", "191", "192", "531", "196", "203", "208", "262", "212", "214", "218", "818", "222", "226", "232", "233", "231", "238", "234", "242", "246",
	"250", "254", "258", "260", "266", "270", "268", "276", "288", "292", "300", "304", "308", "312", "316", "320", "831", "324", "624", "328", "332", "334", "336", "340",
	"348", "352", "356", "360", "364", "368", "372", "833", "376", "380", "388", "392", "832", "400", "398", "404", "296", "408", "410", "414", "417", "418", "428", "422",
	"426", "430", "434", "438", "440", "442", "807", "450", "454", "458", "462", "466", "470", "584", "474", "478", "480", "175", "484", "583", "498", "492", "496", "499",
	"500", "504", "508", "104", "516", "520", "524", "528", "540", "554", "558", "562", "566", "570", "574", "580", "578", "512", "586", "585", "275", "591", "598", "600",
	"604", "608", "612", "616", "620", "630", "634", "638", "642", "643", "646", "652", "654", "659", "662", "663", "666", "670", "882", "674", "678", "682", "686", "688",
	"690", "694", "702", "534", "703", "705", "90", "706", "710", "239", "728", "724", "144", "729", "740", "744", "748", "752", "756", "760", "762", "834", "764", "626",
	"768", "772", "776", "780", "788", "792", "795", "796", "798", "800", "804", "784", "826", "581", "840", "858", "860", "548", "862", "704", "92", "850", "876", "732",
	"887", "894", "716"}

func getRegexFromDictionary(dictionary []string) string {
	return `^(` + strings.Join(dictionary, `|`) + `)$`
}

func regAreaCodeAndContraCode() string {
	return "(" + getRegexFromDictionary(AreaCodeList) + ")|(" + regAreaCode + ")"
}

type CheckColumn struct {
	ColumnName   string
	Description  string
	NotNull      bool
	RegexStr     string
	ColumnType   string
	FuncGetCheck func() []CheckColumn
}

const PACKAGE_NAME = "main."

func CheckSubitem(item interface{}, checkColumns []CheckColumn) map[string]interface{} {
	itemMap := StructToMap(item)
	checkResult := Check(itemMap, checkColumns)
	if checkResult != nil {
		return checkResult
	}
	return nil
}
func Check(dataMap map[string]interface{}, checkColumns []CheckColumn) map[string]interface{} {
	for i := 0; i < len(checkColumns); i++ {
		checkColumn := checkColumns[i]
		item := dataMap[checkColumn.ColumnName]
		fmt.Println(item)
		typeOfA := reflect.TypeOf(item)

		if checkColumn.NotNull {
			if item == nil {
				return buildErrorResult(errors.New(checkColumn.ColumnName + "(" + checkColumn.Description + ")不能为空!"))
			} else if typeOfA.Name() == "string" && item.(string) == "" {
				return buildErrorResult(errors.New(checkColumn.ColumnName + "(" + checkColumn.Description + ")不能为空!"))
			}
		}
		typeString := strings.Replace(typeOfA.String(), PACKAGE_NAME, "", 1)
		// fmt.Println("======")
		// fmt.Println(typeString)
		// fmt.Println(checkColumn.ColumnType)
		// fmt.Println(len(checkColumn.ColumnType))
		if (len(checkColumn.ColumnType) >= 2 && checkColumn.ColumnType[0:2] == "[]") && (len(typeString) < 2 || typeString[0:2] != "[]") {
			return buildErrorResult(errors.New(checkColumn.ColumnName + "(" + checkColumn.Description + ")应当是一个对象数组!"))
		} else if (len(checkColumn.ColumnType) < 2 || checkColumn.ColumnType[0:2] != "[]") && (len(typeString) >= 2 && typeString[0:2] == "[]") {
			return buildErrorResult(errors.New(checkColumn.ColumnName + "(" + checkColumn.Description + ")应当是一个对象而非数组!"))
		} else if typeString == checkColumn.ColumnType {
			// fmt.Println(item)
			if checkColumn.ColumnType[0:2] == "[]" {
				itemArray := reflect.ValueOf(item)
				subCheckColumns := checkColumn.FuncGetCheck()
				for j := 0; j < itemArray.Len(); j++ {
					checkResult := CheckSubitem(itemArray.Index(j).Interface(), subCheckColumns)
					if checkResult != nil {
						return checkResult
					}
				}
			} else {
				checkResult := CheckSubitem(item, checkColumn.FuncGetCheck())
				if checkResult != nil {
					return checkResult
				}
			}
			fmt.Println(item)
		} else if checkColumn.RegexStr != "" && (item != nil && item != 0 && item != "") {
			//fmt.Println(item)
			var matched bool
			if typeOfA.Name() == "string" {
				matched, _ = regexp.MatchString(checkColumn.RegexStr, item.(string))
			} else if typeOfA.Name() == "int" {
				matched, _ = regexp.MatchString(checkColumn.RegexStr, strconv.Itoa(item.(int)))
			}
			if !matched {
				return buildErrorResult(errors.New("参数" + checkColumn.ColumnName + "(" + checkColumn.Description + ")格式不正确!"))
			}
		}
	}
	return nil
}

type Company struct {
	ObjectType          string `json:"docType"`
	CompanyName         string `json:"CompanyName"`
	CompanyCode         string `json:"CompanyCode"`
	RegAdd              string `json:"RegAdd"`
	AreaCode            int    `json:"AreaCode"`
	RegCapital          int    `json:"RegCapital"`
	RegDate             int    `json:"RegDate"`
	BizScope            string `json:"BizScope"`
	TransportNumber     string `json:"TransportNumber"`
	RegTel              string `json:"RegTel"`
	LARName             string `json:"LARName"`
	LARMobile           int    `json:"LARMobile"`
	ContactsName        string `json:"ContactsName"`
	ContactsTel         int    `json:"ContactsTel"`
	EstablishedTime     int    `json:"EstablishedTime"`
	BizLicensePhotoURL  string `json:"BizLicensePhotoURL"`
	ICPNumber           string `json:"ICPNumber"`
	SecurityNumber      string `json:"SecurityNumber"`
	PlatformCompanyName string `json:"PlatformCompanyName"`
	PlatformCompanyCode string `json:"PlatformCompanyCode"`
	Ext                 string `json:"Ext"`
	TransportPhotoURL   string `json:"TransportPhotoURL"`
}

var companyCheck []CheckColumn = nil

func getCompanyCheck() []CheckColumn {
	if companyCheck == nil {
		companyCheck = []CheckColumn{
			CheckColumn{"CompanyName", "公司名称", true, "", "", nil},
			CheckColumn{"CompanyCode", "公司信用代码", true, regCompanyCode, "", nil},
			CheckColumn{"AreaCode", "行政区划代码", true, regAreaCode, "", nil},
			CheckColumn{"PlatformCompanyCode", "平台企业统一社会信用代码", true, "", "", nil},
			CheckColumn{"PlatformCompanyName", "平台企业名称", true, "", "", nil},
			CheckColumn{"RegDate", "营业执照注册日期", false, regAllNumberWithLength(8), "", nil},
			CheckColumn{"EstablishedTime", "成立时间", false, regAllNumberWithLength(8), "", nil},
			CheckColumn{"BizLicensePhotoURL", "营业执照图片", false, regWebsite(), "", nil},
			CheckColumn{"TransportPhotoURL", "道路运输经营许可证图片", false, regWebsite(), "", nil},
		}
	}
	return companyCheck
}

type User struct {
	ObjectType          string `json:"docType"`
	UserID              string `json:"UserID"`
	UserName            string `json:"UserName"`
	UserMobile          string `json:"UserMobile"`
	IDCardFrontPhotoURL string `json:"IDCardFrontPhotoURL"`
	IDCardBackPhotoURL  string `json:"IDCardBackPhotoURL"`
	Ext                 string `json:"Ext"`
}

var userCheck []CheckColumn = nil

func getUserCheck() []CheckColumn {
	if userCheck == nil {
		userCheck = []CheckColumn{
			CheckColumn{"UserID", "身份证或者护照", true, regIdentityOrPassword(), "", nil},
			CheckColumn{"UserName", "用户姓名", true, regChineseOrEnglish(), "", nil},
			CheckColumn{"UserMobile", "手机号", true, regTelPhone(), "", nil},
			CheckColumn{"IDCardFrontPhotoURL", "身份证正面图片", false, regWebsite(), "", nil},
			CheckColumn{"IDCardBackPhotoURL", "身份证背面图片", false, regWebsite(), "", nil},
		}
	}
	return userCheck
}

type Driver struct {
	ObjectType      string `json:"docType"`
	DriverName      string `json:"DriverName"`
	IDNumber        string `json:"IDNumber"`
	VehicleType     string `json:"VehicleType"`
	DLOrg           string `json:"DLOrg"`
	DLValidFrom     int    `json:"DLValidFrom"`
	DLValidUntil    int    `json:"DLValidUntil"`
	QCNumber        string `json:"QCNumber"`
	Mobile          string `json:"Mobile"`
	Notes           string `json:"Notes"`
	DLFrontPhotoURL string `json:"DLFrontPhotoURL"`
	DLBackPhotoURL  string `json:"DLBackPhotoURL"`
	Ext             string `json:"Ext"`
}

var driverCheck []CheckColumn = nil

func getDriverCheck() []CheckColumn {
	if driverCheck == nil {
		driverCheck = []CheckColumn{
			CheckColumn{"IDNumber", "身份证或者护照", true, regIdentityOrPassword(), "", nil},
			CheckColumn{"DriverName", "司机姓名", true, regChineseOrEnglish(), "", nil},
			CheckColumn{"QCNumber", "从业资格证号", true, "", "", nil},
			CheckColumn{"Mobile", "手机号", true, regTelPhone(), "", nil},
			CheckColumn{"VehicleType", "准驾车型", false, regDriverLicense(), "", nil},
			CheckColumn{"DLValidFrom", "驾驶证有效期自", false, regAllNumberWithLength(8), "", nil},
			CheckColumn{"DLValidUntil", "驾驶证有效期至", false, regAllNumberWithLength(8), "", nil},
			CheckColumn{"DLFrontPhotoURL", "驾驶证正面图片", false, regWebsite(), "", nil},
			CheckColumn{"DLBackPhotoURL", "驾驶证副面图片", false, regWebsite(), "", nil},
		}
	}
	return driverCheck
}

type Vehicle struct {
	ObjectType      string `json:"docType"`
	VehicleNumber   string `json:"VehicleNumber"`
	VehicleTypeCode string `json:"VehicleTypeCode"`
	TruckInfo       string `json:"TruckInfo"`
	ShipInfo        string `json:"ShipInfo"`
}

var vehicleCheck []CheckColumn = nil

func getVehicleCheck() []CheckColumn {
	if vehicleCheck == nil {
		vehicleCheck = []CheckColumn{
			CheckColumn{"VehicleNumber", "车辆或船舶号", true, regPlateShip(), "", nil},
			CheckColumn{"VehicleTypeCode", "载具类型代码", true, regVehicleTypeCode, "", nil},
		}
	}
	return vehicleCheck
}

type Order struct {
	ObjectType               string  `json:"docType"`
	Uuid                     string  `json:"Uuid"`
	OrderNumber              string  `json:"OrderNumber"`
	CreateTime               int     `json:"CreateTime"`
	TradeDatetime            int     `json:"TradeDatetime"`
	BizTypeCode              string  `json:"BizTypeCode"`
	OrgTypeCode              int     `json:"OrgTypeCode"`
	PlanDeliveryTimeFrom     int     `json:"PlanDeliveryTimeFrom"`
	PlanDeliveryTimeUntil    int     `json:"PlanDeliveryTimeUntil"`
	PlanArrivalTime          int     `json:"PlanArrivalTime"`
	NeedInvoice              int     `json:"NeedInvoice"`
	PriceType                int     `json:"PriceType"`
	FreightAmount            float64 `json:"FreightAmount"`
	VehicleTypeRequirement   string  `json:"VehicleTypeRequirement"`
	VehicleLengthRequirement string  `json:"VehicleLengthRequirement"`
	GoodsValue               float64 `json:"GoodsValue"`
	Notes                    string  `json:"Notes"`
	Ext                      string  `json:"Ext"`
	CarrierInfo              CarrierInfo
	ConsignorInfo            ConsignorInfo
	SenderInfo               SenderInfo
	ConsigneeInfo            ConsigneeInfo
	GoodsInfo                []GoodsInfo
	Trader                   Trader
	PaymentInfo              []PaymentInfo
	InsuranceInfo            InsuranceInfo
	ContractInfo             ContractInfo
	InvoiceInfo              InvoiceInfo
}

// <Order>
var orderCheck []CheckColumn = nil

func getOrderCheck() []CheckColumn {
	if orderCheck == nil {
		orderCheck = []CheckColumn{
			CheckColumn{"OrderNumber", "委托运输单号", true, "", "", nil},
			CheckColumn{"CreateTime", "生成时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"TradeDatetime", "订单交易时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"BizTypeCode", "业务类型代码", true, regBizTypeCode, "", nil},
			CheckColumn{"OrgTypeCode", "运输组织类型代码", true, regOrgTypeCode, "", nil},
			CheckColumn{"PlanDeliveryTimeFrom", "要求装货起始时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"PlanDeliveryTimeUntil", "要求装货结束时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"PlanArrivalTime", "要求收货时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"PriceType", "价格类型", true, "", "", nil},
			CheckColumn{"FreightAmount", "总运费", true, "", "", nil},
			CheckColumn{"CarrierInfo", "承运人信息", true, "", "CarrierInfo", getCarrierInfoCheck},
			CheckColumn{"ConsignorInfo", "托运人信息", true, "", "ConsignorInfo", getConsignorInfoCheck},
			CheckColumn{"SenderInfo", "发货人信息", true, "", "SenderInfo", getSenderInfoCheck},
			CheckColumn{"ConsigneeInfo", "收货方信息", true, "", "ConsigneeInfo", getConsigneeInfoCheck},
			CheckColumn{"GoodsInfo", "货物信息", true, "", "[]GoodsInfo", getGoodsInfoCheck},
			CheckColumn{"Trader", "交易员信息", true, "", "Trader", getTraderCheck},
			CheckColumn{"PaymentInfo", "付款进度信息", true, "", "[]PaymentInfo", getPaymentInfoCheck},
			CheckColumn{"InsuranceInfo", "保险信息", true, "", "InsuranceInfo", getInsuranceInfoCheck},
			CheckColumn{"ContractInfo", "合同信息", true, "", "ContractInfo", getContractInfoCheck},
			CheckColumn{"InvoiceInfo", "发票信息", true, "", "InvoiceInfo", getInvoiceInfoCheck},
		}
	}
	return orderCheck
}

//承运人信息
type CarrierInfo struct {
	CarrierID         string `json:"CarrierID"`
	CarrierName       string `json:"CarrierName"`
	CarrierBizLicense string `json:"CarrierBizLicense"`
}

var carrierInfoCheck []CheckColumn = nil

func getCarrierInfoCheck() []CheckColumn {
	if carrierInfoCheck == nil {
		carrierInfoCheck = []CheckColumn{
			CheckColumn{"CarrierID", "承运人统一社会信用代码或证件号码", true, "", "", nil},
			CheckColumn{"CarrierName", "承运人名称", true, "", "", nil},
			CheckColumn{"CarrierBizLicense", "承运人道路运输经营许可证编号", true, "", "", nil},
		}
	}
	return carrierInfoCheck
}

//托运人信息
type ConsignorInfo struct {
	ConsignorName string `json:"ConsignorName"`
	ConsignorID   string `json:"ConsignorID"`
}

var consignorInfoCheck []CheckColumn = nil

func getConsignorInfoCheck() []CheckColumn {
	if consignorInfoCheck == nil {
		consignorInfoCheck = []CheckColumn{
			CheckColumn{"ConsignorName", "托运人名称", true, "", "", nil},
			CheckColumn{"ConsignorID", "托运人统一社会信用代码或身份证号", true, "", "", nil},
		}
	}
	return consignorInfoCheck
}

//发货人信息
type SenderInfo struct {
	SenderName      string `json:"SenderName"`
	SenderID        string `json:"SenderID"`
	LoadingAdd      string `json:"LoadingAdd"`
	LoadingAreaCode string `json:"LoadingAreaCode"`
}

var senderInfoCheck []CheckColumn = nil

func getSenderInfoCheck() []CheckColumn {
	if senderInfoCheck == nil {
		senderInfoCheck = []CheckColumn{
			CheckColumn{"SenderName", "发货人名称", true, "", "", nil},
			CheckColumn{"SenderID", "发货人统一社会信用代码或身份证号", true, "", "", nil},
			CheckColumn{"LoadingAdd", "装货地址", true, "", "", nil},
			CheckColumn{"LoadingAreaCode", "装货地点的国家行政区划代码或国别代码", true, regAreaCodeAndContraCode(), "", nil},
		}
	}
	return senderInfoCheck
}

//收货方信息
type ConsigneeInfo struct {
	ConsigneeName   string `json:"ConsigneeName"`
	ConsigneeID     string `json:"ConsigneeID"`
	ArrivalAdd      string `json:"ArrivalAdd"`
	ArrivalAreaCode string `json:"ArrivalAreaCode"`
}

var consigneeInfoCheck []CheckColumn = nil

func getConsigneeInfoCheck() []CheckColumn {
	if consigneeInfoCheck == nil {
		consigneeInfoCheck = []CheckColumn{
			CheckColumn{"ConsigneeName", "收货方名称", true, "", "", nil},
			CheckColumn{"ArrivalAdd", "收货地址", true, "", "", nil},
			CheckColumn{"ArrivalAreaCode", "收货地点的国家行政区划代码或国别代码", true, regAreaCodeAndContraCode(), "", nil},
		}
	}
	return consigneeInfoCheck
}

//货物信息
type GoodsInfo struct {
	GoodsDesc     string  `json:"GoodsDesc"`
	GoodsTypeCode string  `json:"GoodsTypeCode"`
	Weight        float64 `json:"Weight"`
	Cube          float64 `json:"Cube"`
	Quantity      int     `json:"Quantity"`
}

var goodsInfoCheck []CheckColumn = nil

func getGoodsInfoCheck() []CheckColumn {
	if goodsInfoCheck == nil {
		goodsInfoCheck = []CheckColumn{
			CheckColumn{"GoodsDesc", "货物名称", true, "", "", nil},
			CheckColumn{"GoodsTypeCode", "货物类型分类代码", true, regGoodsTypeCode, "", nil},
			CheckColumn{"Weight", "货物名称", true, "", "", nil},
		}
	}
	return goodsInfoCheck
}

//交易员信息
type Trader struct {
	TraderName string `json:"TraderName"`
	TraderID   string `json:"TraderID"`
}

var traderCheck []CheckColumn = nil

func getTraderCheck() []CheckColumn {
	if traderCheck == nil {
		traderCheck = []CheckColumn{
			CheckColumn{"TraderName", "交易员名称", true, "", "", nil},
		}
	}
	return traderCheck
}

//付款进度信息
type PaymentInfo struct {
	SettlementMeansCode string  `json:"SettlementMeansCode"`
	PaymentMeansCode    string  `json:"PaymentMeansCode"`
	PaymentDate         int     `json:"PaymentDate"`
	PaymentAmount       float64 `json:"PaymentAmount"`
}

var paymentInfoCheck []CheckColumn = nil

func getPaymentInfoCheck() []CheckColumn {
	if paymentInfoCheck == nil {
		paymentInfoCheck = []CheckColumn{
			CheckColumn{"SettlementMeansCode", "结算方式代码", true, regSettlementMeansCodeAndObjectTypeCode, "", nil},
			CheckColumn{"PaymentMeansCode", "付款方式代码", true, regPaymentMeansCode, "", nil},
			CheckColumn{"PaymentDate", "付款日期", true, "", "", nil},
			CheckColumn{"PaymentAmount", "付款金额", true, "", "", nil},
		}
	}
	return paymentInfoCheck
}

//保险信息
type InsuranceInfo struct {
	InsuranceNumber      string `json:"InsuranceNumber"`
	InsuranceCompanyCode string `json:"InsuranceCompanyCode"`
}

var insuranceInfoCheck []CheckColumn = nil

func getInsuranceInfoCheck() []CheckColumn {
	if insuranceInfoCheck == nil {
		insuranceInfoCheck = []CheckColumn{
			CheckColumn{"InsuranceNumber", "保险单号", true, "", "", nil},
			CheckColumn{"InsuranceCompanyCode", "保险公司代码", true, getRegexFromDictionary(InsuranceCompanyCodeList), "", nil},
		}
	}
	return insuranceInfoCheck
}

//合同信息
type ContractInfo struct {
	ContractNumber   string `json:"ContractNumber"`
	SigningDatetime  int    `json:"SigningDatetime"`
	ContractPhotoURL string `json:"ContractPhotoURL"`
}

var contractInfoCheck []CheckColumn = nil

func getContractInfoCheck() []CheckColumn {
	if contractInfoCheck == nil {
		contractInfoCheck = []CheckColumn{
			CheckColumn{"ContractNumber", "合同编号", true, "", "", nil},
			CheckColumn{"SigningDatetime", "合同签订时间", true, regAllNumberWithLength(8), "", nil},
			CheckColumn{"ContractPhotoURL", "合同照片地址", true, regWebsite(), "", nil},
		}
	}
	return contractInfoCheck
}

//发票信息
type InvoiceInfo struct {
	InvoiceCode   string `json:"InvoiceCode"`
	InvoiceNumber string `json:"InvoiceNumber"`
}

var invoiceInfoCheck []CheckColumn = nil

func getInvoiceInfoCheck() []CheckColumn {
	if invoiceInfoCheck == nil {
		invoiceInfoCheck = []CheckColumn{
			CheckColumn{"InvoiceCode", "发票代码", false, "", "", nil},
			CheckColumn{"InvoiceNumber", "发票号码", false, "", "", nil},
		}

	}
	return invoiceInfoCheck
}

//合同
type Contract struct {
	ObjectType       string  `json:"docType"`
	Uuid             string  `json:"Uuid"`
	ContractNumber   string  `json:"ContractNumber"`
	ContractType     int     `json:"ContractType"`
	Amount           float64 `json:"Amount"`
	SigningDatetime  int     `json:"SigningDatetime"`
	ContractPhotoURL string  `json:"ContractPhotoURL"`
	Ext              string  `json:"Ext"`
	ConsignorInfo    ConsignorInfo
	CarrierInfo      CarrierInfo
}

var ContractCheck []CheckColumn = nil

func getContractCheck() []CheckColumn {
	if ContractCheck == nil {
		ContractCheck = []CheckColumn{
			CheckColumn{"ContractNumber", "合同单号", true, "", "", nil},
			CheckColumn{"ContractType", "合同类型", true, "", "", nil},
			CheckColumn{"Amount", "合同金额", true, "", "", nil},
			CheckColumn{"SigningDatetime", "合同签订日期", true, regAllNumberWithLength(8), "", nil},
			CheckColumn{"ConsignorInfo", "托运人信息", true, "", "ConsignorInfo", getConsignorInfoCheck},
			CheckColumn{"CarrierInfo", "承运人信息", true, "", "CarrierInfo", getContractCarrierInfoCheck},
			CheckColumn{"ContractPhotoURL", "合同照片", false, regWebsite(), "", nil},
		}

	}
	return ContractCheck
}

//合同的承运人信息check
var ContractCarrierInfoCheck []CheckColumn = nil

func getContractCarrierInfoCheck() []CheckColumn {
	if ContractCarrierInfoCheck == nil {
		ContractCarrierInfoCheck = []CheckColumn{
			CheckColumn{"CarrierID", "承运人统一社会信用代码或证件号码", true, "", "", nil},
			CheckColumn{"CarrierName", "承运人名称", true, "", "", nil},
			CheckColumn{"CarrierBizLicense", "承运人道路运输经营许可证编号", false, "", "", nil},
		}
	}
	return ContractCarrierInfoCheck
}

//发票
type Invoice struct {
	ObjectType      string  `json:"docType"`
	Uuid            string  `json:"Uuid"`
	BuyerName       string  `json:"BuyerName"`
	BuyerTaxNumber  string  `json:"BuyerTaxNumber"`
	BuyerAdd        string  `json:"BuyerAdd"`
	BuyerAccount    string  `json:"BuyerAccount"`
	BuyerTel        string  `json:"BuyerTel"`
	InvoiceDatetime int     `json:"InvoiceDatetime"`
	Drawer          string  `json:"Drawer"`
	SaleName        string  `json:"SaleName"`
	SaleAccount     string  `json:"SaleAccount"`
	SaleTel         string  `json:"SaleTel"`
	SaleAdd         string  `json:"SaleAdd"`
	SaleTaxNumber   string  `json:"SaleTaxNumber"`
	Payee           string  `json:"Payee"`
	Reviewer        string  `json:"Reviewer"`
	InvoiceType     string  `json:"InvoiceType"`
	TotalAmount     float64 `json:"TotalAmount"`
	Notes           float64 `json:"Notes"`
	InvoicePhotoURL string  `json:"InvoicePhotoURL"`
	Ext             string  `json:"Ext"`
	InvoiceDetails  []InvoiceDetails
}

var InvoiceCheck []CheckColumn = nil

func getInvoiceCheck() []CheckColumn {
	if InvoiceCheck == nil {
		InvoiceCheck = []CheckColumn{
			CheckColumn{"BuyerName", "购方名称", true, "", "", nil},
			CheckColumn{"BuyerTaxNumber", "购方税号", true, "", "", nil},
			CheckColumn{"BuyerAdd", "购方地址", true, "", "", nil},
			CheckColumn{"BuyerTel", "购方电话", true, "", "", nil},
			CheckColumn{"SaleName", "销售方", true, "", "", nil},
			CheckColumn{"SaleAccount", "销方银行账号", true, "", "", nil},
			CheckColumn{"SaleTel", "销方电话", true, "", "", nil},
			CheckColumn{"SaleAdd", "销方地址", true, "", "", nil},
			CheckColumn{"SaleTaxNumber", "销方税号", true, "", "", nil},
			CheckColumn{"Payee", "收款人", true, "", "", nil},
			CheckColumn{"Reviewer", "复核人", true, "", "", nil},
			CheckColumn{"InvoiceType", "发票种类", true, "", "", nil},
			CheckColumn{"TotalAmount", "发票金额", true, "", "", nil},
			CheckColumn{"Reviewer", "复核人", true, "", "", nil},
			CheckColumn{"InvoiceDetails", "发票明细信息", true, "", "[]InvoiceDetails", getInvoiceDetailsCheck},
			CheckColumn{"InvoiceDatetime", "开票时间", false, regAllNumberWithLength(14), "", nil},
		}

	}
	return InvoiceCheck
}

//发票明细
type InvoiceDetails struct {
	GoodsName     string  `json:"GoodsName"`
	Amount        int     `json:"Amount"`
	Price         float64 `json:"Price"`
	Unit          string  `json:"Unit"`
	Specification string  `json:"Specification"`
	TaxRate       float64 `json:"TaxRate"`
	TaxAmount     float64 `json:"TaxAmount"`
}

var InvoiceDetailsCheck []CheckColumn = nil

func getInvoiceDetailsCheck() []CheckColumn {
	if InvoiceDetailsCheck == nil {
		InvoiceDetailsCheck = []CheckColumn{
			CheckColumn{"GoodsName", "商品名称", true, "", "", nil},
			CheckColumn{"Amount", "数量", true, "", "", nil},
			CheckColumn{"Price", "单价", true, "", "", nil},
			CheckColumn{"Unit", "单位", true, "", "", nil},
			CheckColumn{"Specification", "规格型号", true, "", "", nil},
			CheckColumn{"TaxRate", "税率", true, "", "", nil},
			CheckColumn{"TaxAmount", "税额", true, "", "", nil},
		}

	}
	return InvoiceDetailsCheck
}

//运单
type Waybill struct {
	ObjectType           string  `json:"docType"`
	Uuid                 string  `json:"Uuid"`
	WaybillNumber        string  `json:"WaybillNumber"`
	OrgTypeCode          int     `json:"OrgTypeCode"`
	FreightAmount        float64 `json:"FreightAmount"`
	BizTypeCode          string  `json:"BizTypeCode"`
	SharingFreightAmount float64 `json:"SharingFreightAmount"`
	Notes                string  `json:"Notes"`
	WaybillCreateTime    int     `json:"WaybillCreateTime"`
	WaybillTradeTime     int     `json:"WaybillTradeTime"`
	Ext                  string  `json:"Ext"`
	DispatchInfo         []DispatchInfo
	DispatcherInfo       DispatcherInfo
	VehicleInfo          VehicleInfo
	DriverInfo           []DriverInfo
	PaymentInfo          []PaymentInfo
	InsuranceInfo        InsuranceInfo
	ContractInfo         ContractInfo
	InvoiceInfo          InvoiceInfo
}

var WaybillCheck []CheckColumn = nil

func getWaybillCheck() []CheckColumn {
	if WaybillCheck == nil {
		WaybillCheck = []CheckColumn{
			CheckColumn{"WaybillNumber", "运单号", true, "", "", nil},
			CheckColumn{"OrgTypeCode", "运输组织类型代码", true, regOrgTypeCode, "", nil},
			CheckColumn{"FreightAmount", "运单金额", true, "", "", nil},
			CheckColumn{"BizTypeCode", "业务类型代码", true, regBizTypeCode, "", nil},
			CheckColumn{"SharingFreightAmount", "分摊运费金额", true, "", "", nil},
			CheckColumn{"WaybillCreateTime", "运单创建时间", true, "", "", nil},
			CheckColumn{"WaybillTradeTime", "运单交易时间", true, "", "", nil},
			CheckColumn{"DispatchInfo", "调度单信息", true, "", "[]DispatchInfo", getDispatchInfoCheck},
			CheckColumn{"DispatcherInfo", "调度员信息", true, "", "DispatcherInfo", getDispatcherInfoCheck},
			CheckColumn{"VehicleInfo", "车辆信息", true, "", "VehicleInfo", getVehicleInfoCheck},
			CheckColumn{"DriverInfo", "驾驶员", true, "", "[]DriverInfo", getDriverInfoCheck},
			CheckColumn{"PaymentInfo", "付款进度信息", true, "", "[]PaymentInfo", getPaymentInfoCheck},
			CheckColumn{"InsuranceInfo", "保险信息", true, "", "InsuranceInfo", getInsuranceInfoCheck},
			CheckColumn{"ContractInfo", "合同信息", true, "", "ContractInfo", getContractInfoCheck},
			CheckColumn{"InvoiceInfo", "发票信息", true, "", "InvoiceInfo", getInvoiceInfoCheck},
		}

	}
	return WaybillCheck
}

//调度单信息
type DispatchInfo struct {
	OrderNumber          string  `json:"OrderNumber"`
	WaybillNumber        string  `json:"WaybillNumber"`
	DispatchNumber       string  `json:"DispatchNumber"`
	PlanDeliveryDatetime int     `json:"PlanDeliveryDatetime"`
	DeliveryDatetime     int     `json:"DeliveryDatetime"`
	PlanArrivalDateTime  int     `json:"PlanArrivalDateTime"`
	ArrivalDateTime      int     `json:"ArrivalDateTime"`
	FreightAmount        float64 `json:"FreightAmount"`
	LoadingPhotoUrl      string  `json:"LoadingPhotoUrl"`
	LadingBillPhotoUrl   string  `json:"LadingBillPhotoUrl"`
	ReceiptPhotoUrl      string  `json:"ReceiptPhotoUrl"`
	Notes                string  `json:"Notes"`
	Ext                  string  `json:"Ext"`
	CarrierInfo          CarrierInfo
	ActualCarrierInfo    ActualCarrierInfo
	ConsignorInfo        ConsignorInfo
	SenderInfo           SenderInfo
	ConsigneeInfo        ConsigneeInfo
	GoodsInfo            []GoodsInfo
}

var DispatchInfoCheck []CheckColumn = nil

func getDispatchInfoCheck() []CheckColumn {
	if DispatchInfoCheck == nil {
		DispatchInfoCheck = []CheckColumn{
			CheckColumn{"OrderNumber", "订单号", true, "", "", nil},
			CheckColumn{"WaybillNumber", "运单号", true, "", "", nil},
			CheckColumn{"DispatchNumber", "调度单号", true, "", "", nil},
			CheckColumn{"PlanDeliveryDatetime", "要求发货日期时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"DeliveryDatetime", "发货日期时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"PlanArrivalDateTime", "税率", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"ArrivalDateTime", "税额", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"FreightAmount", "运费金额", true, "", "", nil},
			CheckColumn{"LoadingPhotoUrl", "装车照片URL", false, regWebsite(), "", nil},
			CheckColumn{"LadingBillPhotoUrl", "提货单照片URL", false, regWebsite(), "", nil},
			CheckColumn{"ReceiptPhotoUrl", "回单照片URL", false, regWebsite(), "", nil},
			CheckColumn{"CarrierInfo", "承运人信息", true, "", "CarrierInfo", getCarrierInfoCheck},
			CheckColumn{"ActualCarrierInfo", "实际承运人信息", true, "", "ActualCarrierInfo", getActualCarrierInfoCheck},
			CheckColumn{"ConsignorInfo", "托运人信息", true, "", "ConsignorInfo", getConsignorInfoCheck},
			CheckColumn{"SenderInfo", "发货人信息", true, "", "SenderInfo", getSenderInfoCheck},
			CheckColumn{"ConsigneeInfo", "收货方信息", true, "", "ConsigneeInfo", getConsigneeInfoCheck},
			CheckColumn{"GoodsInfo", "货物信息", true, "", "[]GoodsInfo", getGoodsInfoCheck},
			CheckColumn{"SenderInfo", "发货人信息", true, "", "SenderInfo", getSenderInfoCheck},
			CheckColumn{"SenderInfo", "发货人信息", true, "", "SenderInfo", getSenderInfoCheck},
		}

	}
	return DispatchInfoCheck
}

//实际承运人信息
type ActualCarrierInfo struct {
	ActualCarrierName       string `json:"ActualCarrierName"`
	ActualCarrierBizLicense string `json:"ActualCarrierBizLicense"`
	ActualCarrierID         string `json:"ActualCarrierID"`
}

var ActualCarrierInfoCheck []CheckColumn = nil

func getActualCarrierInfoCheck() []CheckColumn {
	if ActualCarrierInfoCheck == nil {
		ActualCarrierInfoCheck = []CheckColumn{
			CheckColumn{"ActualCarrierID", "实际承运人统一社会信用代码或证件号码", true, "", "", nil},
			CheckColumn{"ActualCarrierName", "实际承运人名称", true, "", "", nil},
			CheckColumn{"ActualCarrierBizLicense", "实际承运人道路运输经营许可证号", true, "", "", nil},
		}

	}
	return ActualCarrierInfoCheck
}

//调度员信息
type DispatcherInfo struct {
	DispaccherName string `json:"DispaccherName"`
	DispatcherID   string `json:"DispatcherID"`
}

var DispatcherInfoCheck []CheckColumn = nil

func getDispatcherInfoCheck() []CheckColumn {
	if DispatcherInfoCheck == nil {
		DispatcherInfoCheck = []CheckColumn{
			CheckColumn{"DispaccherName", "调度员名称", true, "", "", nil},
			CheckColumn{"DispatcherID", "调度员身份证号", false, regIdentityOrPassword(), "", nil},
		}

	}
	return DispatcherInfoCheck
}

//车辆信息
type VehicleInfo struct {
	VehicleNumber   string `json:"VehicleNumber"`
	VehicleTypeCode string `json:"VehicleTypeCode"`
}

var VehicleInfoCheck []CheckColumn = nil

func getVehicleInfoCheck() []CheckColumn {
	if VehicleInfoCheck == nil {
		VehicleInfoCheck = []CheckColumn{
			CheckColumn{"VehicleNumber", "车辆或船舶号", true, "", "", nil},
			CheckColumn{"VehicleTypeCode", "载具类型代码", true, regVehicleTypeCode, "", nil},
		}

	}
	return VehicleInfoCheck
}

//驾驶员
type DriverInfo struct {
	DriverName string `json:"DriverName"`
	IDNumber   string `json:"IDNumber"`
}

var DriverInfoCheck []CheckColumn = nil

func getDriverInfoCheck() []CheckColumn {
	if DriverInfoCheck == nil {
		DriverInfoCheck = []CheckColumn{
			CheckColumn{"DriverName", "姓名", true, "", "", nil},
			CheckColumn{"IDNumber", "身份证号", true, regIdentityOrPassword(), "", nil},
		}

	}
	return DriverInfoCheck
}

//结算
type CapitalFlow struct {
	ObjectType            string  `json:"docType"`
	Uuid                  string  `json:"Uuid"`
	ARAPNumber            string  `json:"ARAPNumber"`
	CreateTime            int     `json:"CreateTime"`
	VehicleNumber         string  `json:"VehicleNumber"`
	VehicleTypeCode       string  `json:"VehicleTypeCode"`
	VehiclePlateColorCode string  `json:"VehiclePlateColorCode"`
	WaybillTypeCode       int     `json:"WaybillTypeCode"`
	WaybillNumber         string  `json:"WaybillNumber"`
	SettlementMeansCode   string  `json:"SettlementMeansCode"`
	PaymentMeansCode      string  `json:"PaymentMeansCode"`
	PaymentDate           int     `json:"PaymentDate"`
	PaymentAmount         float64 `json:"PaymentAmount"`
	Notes                 string  `json:"Notes"`
	Ext                   string  `json:"Ext"`
	Payer                 Payer
	Payee                 Payee
	CapitalflowList       []CapitalflowList
}

var CapitalFlowCheck []CheckColumn = nil

func getCapitalFlowCheck() []CheckColumn {
	if CapitalFlowCheck == nil {
		CapitalFlowCheck = []CheckColumn{
			CheckColumn{"ARAPNumber", "应结单据号", true, "", "", nil},
			CheckColumn{"CreateTime", "创建时间", true, regAllNumberWithLength(14), "", nil},
			CheckColumn{"VehicleNumber", "车辆或船舶号", true, "", "", nil},
			CheckColumn{"VehicleTypeCode", "载具类型代码", true, regVehicleTypeCode, "", nil},
			CheckColumn{"VehiclePlateColorCode", "车牌颜色代码", false, regPlateColorCode, "", nil},
			CheckColumn{"WaybillTypeCode", "运输单类型代码", true, regWaybillTypeCode, "", nil},
			CheckColumn{"WaybillNumber", "托运单号", true, "", "", nil},
			CheckColumn{"SettlementMeansCode", "结算方式代码", true, regSettlementMeansCodeAndObjectTypeCode, "", nil},
			CheckColumn{"PaymentMeansCode", "付款方式代码", true, regPaymentMeansCode, "", nil},
			CheckColumn{"PaymentDate", "应付日期", true, regAllNumberWithLength(8), "", nil},
			CheckColumn{"PaymentAmount", "应付金额", true, "", "", nil},
			CheckColumn{"Payer", "付款方公司", true, "", "Payer", getPayerCheck},
			CheckColumn{"Payee", "收款方公司", true, "", "Payee", getPayeeCheck},
			CheckColumn{"CapitalflowList", "付款方公司", true, "", "[]CapitalflowList", getCapitalflowListCheck},
		}

	}
	return CapitalFlowCheck
}

//付款方公司
type Payer struct {
	PayerID   string `json:"PayerID"`
	PayerName string `json:"PayerName"`
}

var PayerCheck []CheckColumn = nil

func getPayerCheck() []CheckColumn {
	if PayerCheck == nil {
		PayerCheck = []CheckColumn{
			CheckColumn{"PayerID", "付款方统一社会信用代码或证件号码", true, "", "", nil},
			CheckColumn{"PayerName", "收款方名称", true, "", "", nil},
		}

	}
	return PayerCheck
}

//收款方公司
type Payee struct {
	PayeeID   string `json:"PayeeID"`
	PayeeName string `json:"PayeeName"`
}

var PayeeCheck []CheckColumn = nil

func getPayeeCheck() []CheckColumn {
	if PayeeCheck == nil {
		PayeeCheck = []CheckColumn{
			CheckColumn{"PayeeID", "收款方统一社会信用代码或证件号", true, "", "", nil},
			CheckColumn{"PayeeName", "收款方名称", true, "", "", nil},
		}

	}
	return PayeeCheck
}

//资金流水列表
type CapitalflowList struct {
	PaymentMeansCode   string  `json:"PaymentMeansCode"`
	ActualPayeeName    string  `json:"ActualPayeeName"`
	ActualPayeeAccount string  `json:"ActualPayeeAccount"`
	PayeeBankCode      string  `json:"PayeeBankCode"`
	SerialNumber       string  `json:"SerialNumber"`
	PaymentAmount      float64 `json:"PaymentAmount"`
	PaymentTime        int     `json:"PaymentTime"`
}

var CapitalflowListCheck []CheckColumn = nil

func getCapitalflowListCheck() []CheckColumn {
	if CapitalflowListCheck == nil {
		CapitalflowListCheck = []CheckColumn{
			CheckColumn{"PaymentMeansCode", "付款方式代码", true, regPaymentMeansCode, "", nil},
			CheckColumn{"ActualPayeeName", "实际收款方名称", true, "", "", nil},
			CheckColumn{"PayeeBankCode", "实际收款方银行代码", false, getRegexFromDictionary(BankCodeList), "", nil},
			CheckColumn{"ActualPayeeAccount", "实际收款帐户信息", true, "", "", nil},
			CheckColumn{"SerialNumber", "资金流水号/序列号", true, "", "", nil},
			CheckColumn{"PaymentAmount", "实际支付金额", true, "", "", nil},
			CheckColumn{"PaymentTime", "日期时间", true, regAllNumberWithLength(14), "", nil},
		}

	}
	return CapitalflowListCheck
}

//投诉
type Complaint struct {
	ObjectType     string `json:"docType"`
	Uuid           string `json:"Uuid"`
	ObjectTypeCode int    `json:"ObjectTypeCode"`
	ObjectID       string `json:"ObjectID"`
	Content        string `json:"Content"`
	CreateTime     int    `json:"CreateTime"`
	Result         string `json:"Result"`
	Ext            string `json:"Ext"`
}

var ComplaintCheck []CheckColumn = nil

func getComplaintCheck() []CheckColumn {
	if ComplaintCheck == nil {
		ComplaintCheck = []CheckColumn{
			CheckColumn{"ObjectTypeCode", "投诉对象类别代码", true, regSettlementMeansCodeAndObjectTypeCode, "", nil},
			CheckColumn{"ObjectID", "投诉对象ID", true, "", "", nil},
			CheckColumn{"Content", "投诉内容", true, "", "", nil},
			CheckColumn{"CreateTime", "投诉时间", true, "", regAllNumberWithLength(14), nil},
			CheckColumn{"Result", "投诉处理结果", true, "", "", nil},
		}

	}
	return ComplaintCheck
}

//评价
type Comment struct {
	ObjectType     string `json:"docType"`
	Uuid           string `json:"Uuid"`
	ObjectTypeCode int    `json:"ObjectTypeCode"`
	ObjectID       string `json:"ObjectID"`
	Level          int    `json:"Level"`
	ContentType    int    `json:"ContentType"`
	Content        string `json:"Content"`
	Valuator       string `json:"Valuator"`
	CreateTime     string `json:"CreateTime"`
	Notes          string `json:"Notes"`
	Ext            string `json:"Ext"`
}

var CommentCheck []CheckColumn = nil

func getCommentCheck() []CheckColumn {
	if CommentCheck == nil {
		CommentCheck = []CheckColumn{
			CheckColumn{"ObjectTypeCode", "评价对象类别", true, regSettlementMeansCodeAndObjectTypeCode, "", nil},
			CheckColumn{"ObjectID", "评价对象ID", true, "", "", nil},
			CheckColumn{"Level", "评价等级代码", true, regLevel, "", nil},
			CheckColumn{"ContentType", "评价内容分类", true, "", "", nil},
			CheckColumn{"Valuator", "评价人", true, "", "", nil},
			CheckColumn{"CreateTime", "评价时间", true, regAllNumberWithLength(12), "", nil},
		}

	}
	return CommentCheck
}

//company添加
func (sc *SimpleContract) AddCompany(ctx contractapi.TransactionContextInterface, company string) (map[string]interface{}, error) {
	var com Company
	json.Unmarshal([]byte(company), &com) //string转struct
	com.ObjectType = "company"
	companyData := StructToMap(com) //struct转map
	checkResult := Check(companyData, getCompanyCheck())
	if checkResult != nil {
		return checkResult, nil
	}
	companyCode := companyData["CompanyCode"]
	existing, err := ctx.GetStub().GetState(companyCode.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	companyAsBytes, err := json.Marshal(com)
	err = ctx.GetStub().PutState(companyCode.(string), companyAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//company更新
func (sc *SimpleContract) ModifyCompany(ctx contractapi.TransactionContextInterface, company string) (map[string]interface{}, error) {

	var com Company
	json.Unmarshal([]byte(company), &com) //string转struct
	com.ObjectType = "company"
	companyData := StructToMap(com) //struct转map
	checkResult := Check(companyData, getCompanyCheck())
	if checkResult != nil {
		return checkResult, nil
	}

	companyCode := companyData["CompanyCode"]
	bytes, err := ctx.GetStub().GetState(companyCode.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if bytes == nil {
		return buildErrorResult(errors.New("修改的数据不存在!")), nil
	}
	comAsBytes, err := json.Marshal(com)
	err = ctx.GetStub().PutState(companyCode.(string), comAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//company富查询
func (sc *SimpleContract) QueryCompany(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "company")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//user添加
func (sc *SimpleContract) AddUser(ctx contractapi.TransactionContextInterface, users string) (map[string]interface{}, error) {

	var user User
	json.Unmarshal([]byte(users), &user) //string转struct
	user.ObjectType = "user"
	userMapData := StructToMap(user) //struct转map
	checkResult := Check(userMapData, getUserCheck())
	if checkResult != nil {
		return checkResult, nil
	}

	userID := userMapData["UserID"]
	existing, err := ctx.GetStub().GetState(userID.(string) + "user")
	if err != nil {
		return buildErrorResult(err), nil
	}

	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}

	userAsBytes, err := json.Marshal(user)
	err = ctx.GetStub().PutState(userID.(string)+"user", userAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//user更新
func (sc *SimpleContract) ModifyUser(ctx contractapi.TransactionContextInterface, cond string) (map[string]interface{}, error) {

	var user User
	json.Unmarshal([]byte(cond), &user) //string转struct
	user.ObjectType = "user"
	userMapData := StructToMap(user) //struct转map
	checkResult := Check(userMapData, getUserCheck())
	if checkResult != nil {
		return checkResult, nil
	}
	userID := userMapData["UserID"]
	bytes, err := ctx.GetStub().GetState(userID.(string) + "user")
	if err != nil {
		return buildErrorResult(err), nil
	}
	if bytes == nil {
		return buildErrorResult(errors.New("修改的数据不存在!")), nil
	}

	userAsBytes, err := json.Marshal(user)
	err = ctx.GetStub().PutState(userID.(string)+"user", userAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//user富查询
func (sc *SimpleContract) QueryUser(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "user")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Driver添加

func (sc *SimpleContract) AddDriver(ctx contractapi.TransactionContextInterface, drivers string) (map[string]interface{}, error) {
	var dirver Driver
	json.Unmarshal([]byte(drivers), &dirver) //string转struct
	dirver.ObjectType = "dirver"
	dirverData := StructToMap(dirver) //struct转map
	checkResult := Check(dirverData, getDriverCheck())
	if checkResult != nil {
		return checkResult, nil
	}
	dirverNumber := dirverData["IDNumber"]
	existing, err := ctx.GetStub().GetState(dirverNumber.(string) + "driver")
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	dirverAsBytes, err := json.Marshal(dirver)
	err = ctx.GetStub().PutState(dirverNumber.(string)+"driver", dirverAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//Driver更新
func (sc *SimpleContract) ModifyDriver(ctx contractapi.TransactionContextInterface, dirvers string) (map[string]interface{}, error) {

	var dirver Driver
	json.Unmarshal([]byte(dirvers), &dirver) //string转struct
	dirver.ObjectType = "dirver"
	dirverMapData := StructToMap(dirver) //struct转map
	checkResult := Check(dirverMapData, getDriverCheck())
	if checkResult != nil {
		return checkResult, nil
	}

	dirverNumber := dirverMapData["IDNumber"]
	bytes, err := ctx.GetStub().GetState(dirverNumber.(string) + "driver")
	if err != nil {
		return buildErrorResult(err), nil
	}
	if bytes == nil {
		return buildErrorResult(errors.New("修改的数据不存在!")), nil
	}

	dirverAsBytes, err := json.Marshal(dirver)
	err = ctx.GetStub().PutState(dirverNumber.(string)+"driver", dirverAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//Driver富查询
func (sc *SimpleContract) QueryDriver(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "dirver")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Vehicle添加
func (sc *SimpleContract) AddVehicle(ctx contractapi.TransactionContextInterface, vehicles string) (map[string]interface{}, error) {

	var vehicle Vehicle
	json.Unmarshal([]byte(vehicles), &vehicle) //string转struct
	vehicle.ObjectType = "vehicle"
	vehicleMapData := StructToMap(vehicle) //struct转map
	checkResult := Check(vehicleMapData, getVehicleCheck())

	if checkResult != nil {
		return checkResult, nil
	}
	vehicleNumber := vehicleMapData["VehicleNumber"]
	existing, err := ctx.GetStub().GetState(vehicleNumber.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}

	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}

	vehicleAsBytes, err := json.Marshal(vehicle)
	err = ctx.GetStub().PutState(vehicleNumber.(string), vehicleAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//Vehicle更新
func (sc *SimpleContract) ModifyVehicle(ctx contractapi.TransactionContextInterface, vehicles string) (map[string]interface{}, error) {

	var vehicle Vehicle
	json.Unmarshal([]byte(vehicles), &vehicle) //string转struct
	vehicle.ObjectType = "vehicle"
	vehicleMapData := StructToMap(vehicle) //struct转map
	checkResult := Check(vehicleMapData, getVehicleCheck())
	if checkResult != nil {
		return checkResult, nil
	}
	vehicleNumber := vehicleMapData["VehicleNumber"]
	bytes, err := ctx.GetStub().GetState(vehicleNumber.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if bytes == nil {
		return buildErrorResult(errors.New("修改的数据不存在!")), nil
	}
	vehicleAsBytes, err := json.Marshal(vehicle)
	err = ctx.GetStub().PutState(vehicleNumber.(string), vehicleAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	result["code"] = 200
	return result, nil
}

//Vehicle富查询
func (sc *SimpleContract) QueryVehicle(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "vehicle")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Order新增
func (sc *SimpleContract) AddOrder(ctx contractapi.TransactionContextInterface, orders string) (map[string]interface{}, error) {

	var order Order
	json.Unmarshal([]byte(orders), &order) //string转struct
	for i := 0; i < len(order.GoodsInfo); i++ {
		weight := order.GoodsInfo[i].Weight
		wei := fmt.Sprintf("%.3f", weight)
		weightFloat, err := strconv.ParseFloat(wei, 64)
		if err != nil {
			return buildErrorResult(err), nil
		}
		order.GoodsInfo[i].Weight = weightFloat
		cube := order.GoodsInfo[i].Cube
		cub := fmt.Sprintf("%.4f", cube)
		cubeFloat, err := strconv.ParseFloat(cub, 64)
		if err != nil {
			return buildErrorResult(err), nil
		}
		order.GoodsInfo[i].Cube = cubeFloat
	}
	for i := 0; i < len(order.PaymentInfo); i++ {
		paymentAmount := order.PaymentInfo[i].PaymentAmount
		pay := fmt.Sprintf("%.3f", paymentAmount)
		payfloat, err := strconv.ParseFloat(pay, 64)
		if err != nil {
			return buildErrorResult(err), nil
		}
		order.PaymentInfo[i].PaymentAmount = payfloat
	}
	order.ObjectType = "order"
	orderData := StructToMap(order) //struct转map
	orderCheckResult := Check(orderData, getOrderCheck())
	if orderCheckResult != nil {
		return orderCheckResult, nil
	}
	uuid := orderData["Uuid"]
	existing, err := ctx.GetStub().GetState(uuid.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(uuid.(string), orderAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = uuid.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Order富查询
func (sc *SimpleContract) QueryOrder(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "order")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Contract新增
func (sc *SimpleContract) AddContract(ctx contractapi.TransactionContextInterface, contracts string) (map[string]interface{}, error) {

	var contract Contract
	json.Unmarshal([]byte(contracts), &contract) //string转struct
	contract.ObjectType = "contract"
	amount := contract.Amount
	mount := fmt.Sprintf("%.3f", amount)
	mountFloat, err := strconv.ParseFloat(mount, 64)
	if err != nil {
		return buildErrorResult(err), nil
	}
	contract.Amount = mountFloat
	contractData := StructToMap(contract) //struct转map
	orderCheckResult := Check(contractData, getContractCheck())
	if orderCheckResult != nil {
		return orderCheckResult, nil
	}

	uuid := contractData["Uuid"]
	existing, err := ctx.GetStub().GetState(uuid.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	contractAsBytes, err := json.Marshal(contract)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(uuid.(string), contractAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = uuid.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Contract富查询
func (sc *SimpleContract) QueryContract(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "contract")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Invoice新增
func (sc *SimpleContract) AddInvoice(ctx contractapi.TransactionContextInterface, invoices string) (map[string]interface{}, error) {

	var invoice Invoice
	json.Unmarshal([]byte(invoices), &invoice) //string转struct
	invoice.ObjectType = "invoice"
	for i := 0; i < len(invoice.InvoiceDetails); i++ {
		price := invoice.InvoiceDetails[i].Price
		pri := fmt.Sprintf("%.3f", price)
		priceFloat, err := strconv.ParseFloat(pri, 64)
		if err != nil {
			return buildErrorResult(err), nil
		}
		invoice.InvoiceDetails[i].Price = priceFloat
		taxAmount := invoice.InvoiceDetails[i].TaxAmount
		mount := fmt.Sprintf("%.3f", taxAmount)
		mountFloat, err := strconv.ParseFloat(mount, 64)
		if err != nil {
			return buildErrorResult(err), nil
		}
		invoice.InvoiceDetails[i].TaxAmount = mountFloat

	}
	invoiceData := StructToMap(invoice) //struct转map
	invoiceCheckResult := Check(invoiceData, getInvoiceCheck())
	if invoiceCheckResult != nil {
		return invoiceCheckResult, nil
	}

	uuid := invoiceData["Uuid"]
	existing, err := ctx.GetStub().GetState(uuid.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	invoiceAsBytes, err := json.Marshal(invoice)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(uuid.(string), invoiceAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = uuid.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Invoice富查询
func (sc *SimpleContract) QueryInvoice(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "invoice")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Waybill运单新增
func (sc *SimpleContract) AddWaybill(ctx contractapi.TransactionContextInterface, Waybills string) (map[string]interface{}, error) {

	var waybill Waybill
	json.Unmarshal([]byte(Waybills), &waybill) //string转struct
	waybill.ObjectType = "waybill"
	price := waybill.SharingFreightAmount
	pri := fmt.Sprintf("%.3f", price)
	priceFloat, err := strconv.ParseFloat(pri, 64)
	if err != nil {
		return buildErrorResult(err), nil
	}
	waybill.SharingFreightAmount = priceFloat
	waybillData := StructToMap(waybill) //struct转map
	waybillCheckResult := Check(waybillData, getWaybillCheck())
	if waybillCheckResult != nil {
		return waybillCheckResult, nil
	}

	uuid := waybillData["Uuid"]
	existing, err := ctx.GetStub().GetState(uuid.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	waybillAsBytes, err := json.Marshal(waybill)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(uuid.(string), waybillAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = uuid.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Waybill运单富查询
func (sc *SimpleContract) QueryWaybill(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "waybill")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//CapitalFlow结算新增
func (sc *SimpleContract) AddCapitalFlow(ctx contractapi.TransactionContextInterface, capitalFlows string) (map[string]interface{}, error) {
	var capitalFlow CapitalFlow
	json.Unmarshal([]byte(capitalFlows), &capitalFlow) //string转struct
	capitalFlow.ObjectType = "capitalFlow"

	capitalFlowData := StructToMap(capitalFlow) //struct转map
	capitalFlowCheckResult := Check(capitalFlowData, getCapitalFlowCheck())
	if capitalFlowCheckResult != nil {
		return capitalFlowCheckResult, nil
	}

	uuid := capitalFlowData["Uuid"]
	existing, err := ctx.GetStub().GetState(uuid.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	capitalFlowAsBytes, err := json.Marshal(capitalFlow)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(uuid.(string), capitalFlowAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = uuid.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil

}

//CapitalFlow结算查询
func (sc *SimpleContract) QueryCapitalFlow(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "capitalFlow")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Complaint投诉新增
func (sc *SimpleContract) AddComplaint(ctx contractapi.TransactionContextInterface, complaints string) (map[string]interface{}, error) {

	var complaint Complaint
	json.Unmarshal([]byte(complaints), &complaint) //string转struct
	complaint.ObjectType = "complaint"
	complaintData := StructToMap(complaint) //struct转map
	complaintCheckResult := Check(complaintData, getComplaintCheck())
	if complaintCheckResult != nil {
		return complaintCheckResult, nil
	}

	id := complaintData["Uuid"]
	existing, err := ctx.GetStub().GetState(id.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	complaintAsBytes, err := json.Marshal(complaint)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(id.(string), complaintAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = id.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Complaint投诉查询
func (sc *SimpleContract) QueryComplaint(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "complaint")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

//Comment评价新增
func (sc *SimpleContract) AddComment(ctx contractapi.TransactionContextInterface, comments string) (map[string]interface{}, error) {

	var comment Comment
	json.Unmarshal([]byte(comments), &comment) //string转struct
	comment.ObjectType = "comment"
	commentData := StructToMap(comment) //struct转map
	commentCheckResult := Check(commentData, getCommentCheck())
	if commentCheckResult != nil {
		return commentCheckResult, nil
	}

	id := commentData["Uuid"]
	existing, err := ctx.GetStub().GetState(id.(string))
	if err != nil {
		return buildErrorResult(err), nil
	}
	if existing != nil {
		return buildErrorResult(errors.New("添加的数据已经存在!")), nil
	}
	commentAsBytes, err := json.Marshal(comment)
	if err != nil {
		return buildErrorResult(err), nil
	}
	err = ctx.GetStub().PutState(id.(string), commentAsBytes)
	if err != nil {
		return buildErrorResult(err), nil
	}
	result := make(map[string]interface{})
	dataUuidMap := make(map[string]interface{})
	dataUuidMap["Uuid"] = id.(string)
	result["code"] = 200
	result["data"] = dataUuidMap
	return result, nil
}

//Comment评价查询
func (sc *SimpleContract) QueryComment(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	queryString, err := buildSelector(cond, "comment")
	if err != nil {
		return buildErrorResult(err), nil
	}
	page, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return buildErrorResult(err), nil
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(page), bookmark) // 富查询的返回结果可能为多条 所以这里返回的是一个迭代器 需要我们进一步的处理来获取需要的结果
	if err != nil {
		return buildErrorResult(err), nil
	}
	defer resultsIterator.Close() //释放迭代器

	entityList, err := buidMapListFromIterator(resultsIterator)
	if err != nil {
		return buildErrorResult(err), nil
	}

	result := make(map[string]interface{})
	result["data"] = entityList
	result["responseMetadata"] = responseMetadata
	result["code"] = 200
	return result, nil
}

// daniu v2 - company富查询
func (sc *SimpleContract) QueryCompanyV2(ctx contractapi.TransactionContextInterface, cond string, pageSize string, bookmark string) (map[string]interface{}, error) {
	return QueryCompany(ctx, cond, pageSize, bookmark)
}


// func test(capitalFlows string)(map[string]interface{}, error){

//     var capitalFlow CapitalFlow
//     json.Unmarshal([]byte(capitalFlows), &capitalFlow) //string转struct
//     capitalFlow.ObjectType = "capitalFlow"
//     capitalFlowData := StructToMap(capitalFlow)      //struct转map
//     capitalFlowCheckResult := Check(capitalFlowData, getCapitalFlowCheck())
//     fmt.Println(capitalFlowCheckResult)
//     if capitalFlowCheckResult != nil {
//         return capitalFlowCheckResult, nil
//     }

//     result := make(map[string]interface{})

//     result["code"] = 200

//     return result, nil
// }

func main() {

	// var a float64
	// a = 12
	// s := fmt.Sprintf("%.4f",a)
	// fmt.Println(s)
	// typeOfA := reflect.TypeOf(s)
	// fmt.Println(typeOfA.Name())
	// float,err := strconv.ParseFloat(s, 64)
	// if err != nil{
	// }
	//fmt.Println(float)
	// matched, _ := regexp.MatchString(`^[A-Z0-9]{18}$`, "14103419960S270012")
	// fmt.Println(matched)
	// var str string
	// str = "{\"ARAPNumber\":\"20191213028\",\"CreateTime\":20191213110201,\"Payer\":{\"PayerID\":\"919432432432432432F\",\"PayerName\":\"付款人公司\"},\"Payee\":{\"PayeeID\":\"919432432432432432F\",\"PayeeName\":\"收款人公司\"},\"VehicleNumber\":\"苏AH3295\",\"VehicleTypeCode\":\"2\",\"WaybillType\":1,\"WaybillNumber\":\"20202222221\",\"SettlementMeansCode\":\"12\",\"PaymentMeansCode\":\"42\",\"PaymentDate\":20200420,\"PaymentAmount\":999.999,\"CapitalflowList\":[{\"PaymentMeansCode\":\"42\",\"ActualPayeeName\":\"黄药师\",\"ActualPayeeAccount\":\"13017603285\",\"PayeeBankCode\":\"BKCH\",\"SerialNumber\":\"CW191213028\",\"PaymentAmount\":30000.000,\"PaymentTime\":20190926084701}],\"Notes\":\"\",\"Ext\":\"扩展信息\"}"
	// test(str)

	contract := new(SimpleContract)
	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		panic("创建智能合约失败：" + err.Error())
	}
	if err := cc.Start(); err != nil {
		panic("启动智能合约失败：" + err.Error())
	}

}

func (sc *SimpleContract) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("SimpleContract Init")
	return nil
}
