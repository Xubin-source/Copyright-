package main

import (
	//"bytes"
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	ip "github.com/hyperledger/fabric/protos/peer"
)

type IntelProChainCode struct {
}

//information of the production
type AppInfo struct {
	Applicantname string `json:appname` //name of applicant //也是法定归属人
	Applicantid   string `json:appid`   //id number of applicant
	Usingprice    string `json:usingprice`
	Purprice      string `json:purprice` //price of purchasing the production
	//Appaccount string `json:appacount` //account balance of applicant
	Contribution string `json:appper` //the percentage of contributin of applicant

	PartnerInfo []PartInfo `json:partnerInfo` //information of partner 合作人信息，结构体数组
	Hashvalue   string     `json:hashvalue`
}

//information of partner
type PartInfo struct {
	Name string `json:parname`
	Id   string `json:parid`
	Con  string `json:parper` // the percentage of contributin of partner合作者贡献占比

}

type QueryResult struct {
	ContentNum string `json:"ContentNum"`
	Record     *AppInfo
}

type Person struct {
	Name    string `json:name`
	Id      string `json:id`
	Account string `json:account`
}

func (t *IntelProChainCode) Invoke(stub shim.ChaincodeStubInterface) ip.Response {
	fun, args := stub.GetFunctionAndParameters()

	if fun == "queryuserinfo" {
		return t.queryuseinfo(stub, args)
	} else if fun == "queryallproinfo" {
		return t.queryallproinfo(stub)
	} else if fun == "querycontinfo" {
		return t.querycontinfo(stub)
	} else if fun == "submit" {
		return t.submit(stub, args)
	} else if fun == "subscribe" {
		return t.subscribe(stub, args)
	} else if fun == "buy" {
		return t.buy(stub, args)
	} else if fun == "addPerInfo" {
		return t.addPerInfo(stub, args)
	} else if fun == "queryAllPerInfo" {
		return t.queryAllPerInfo(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}

func (t *IntelProChainCode) Init(APIstub shim.ChaincodeStubInterface) ip.Response {
	production := []AppInfo{
		AppInfo{Applicantname: "A", Applicantid: "000", Usingprice: "20", Purprice: "500", Contribution: "50", PartnerInfo: []PartInfo{PartInfo{Name: "A1", Id: "001", Con: "50"}}, Hashvalue: "XXX"},
		AppInfo{Applicantname: "B", Applicantid: "100", Usingprice: "50", Purprice: "", Contribution: "50", PartnerInfo: []PartInfo{PartInfo{Name: "B1", Id: "101", Con: "30"}, PartInfo{Name: "B2", Id: "102", Con: "20"}}, Hashvalue: "XXX"},
	}
	person := []Person{
		Person{Name: "A", Id: "000", Account: "100"},
		Person{Name: "A1", Id: "001", Account: "100"},
		Person{Name: "B", Id: "100", Account: "50"},
		Person{Name: "B1", Id: "101", Account: "100"},
		Person{Name: "B2", Id: "102", Account: "70"},
		Person{Name: "C", Id: "200", Account: "150"},
		Person{Name: "D", Id: "300", Account: "1000"},
		Person{Name: "E", Id: "400", Account: "1000"},
	}
	for i, pro := range production {
		appinfoBytes, _ := json.Marshal(pro)
		APIstub.PutState("cont"+strconv.Itoa(i), appinfoBytes)

	}
	for _, per := range person {
		personinfo, _ := json.Marshal(per)
		APIstub.PutState(per.Id, personinfo)
	}
	return shim.Success(nil)

}

//add person information
func (t *IntelProChainCode) addPerInfo(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {
	if len(args) != 3 {
		return shim.Error("Please input correct information of person")
	}
	var perinfo = Person{Name: args[1], Id: args[2], Account: args[3]}
	perinfores, _ := json.Marshal(perinfo)

	APIstub.PutState(args[0], perinfores)

	return shim.Success(nil)
}

//apply information about applicant 添加申请信息
func (t *IntelProChainCode) submit(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {
	if (len(args)-7)%3 != 0 {
		return shim.Error("Please input the correct format (pro number,your information)")
	}

	var parinfo PartInfo
	var appinfo = AppInfo{Applicantname: args[1], Applicantid: args[2], Usingprice: args[3], Purprice: args[4], Contribution: args[5], Hashvalue: args[9]}

	if len(args) == 7 {

		appinfores, _ := json.Marshal(appinfo)
		APIstub.PutState(args[0], appinfores)
		return shim.Success(nil)
	}

	for i := 6; i < len(args)-1; {
		parinfo.Name = args[i]
		parinfo.Id = args[i+1]
		parinfo.Con = args[i+2]

		appinfo.PartnerInfo = append(appinfo.PartnerInfo, parinfo)
		i = i + 3
	}
	appinfores, err := json.Marshal(appinfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	APIstub.PutState(args[0], appinfores)
	return shim.Success(nil)
}

//query all information of person
func (t *IntelProChainCode) queryAllPerInfo(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {
	startKey := "000"
	endKey := "999"
	var perinfo Person
	results := []Person{}
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}
		_ = json.Unmarshal(queryResponse.Value, &perinfo)
		results = append(results, perinfo)
	}
	perAsBytes, err := json.Marshal(results)
	return shim.Success(perAsBytes)
}

//query all information of production in the ledger 查询所有产品信息
func (t *IntelProChainCode) queryallproinfo(APIstub shim.ChaincodeStubInterface) ip.Response {
	startKey := "cont0"
	endKey := "cont99"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}

		pro := new(AppInfo)
		_ = json.Unmarshal(queryResponse.Value, pro)

		queryResult := QueryResult{ContentNum: queryResponse.Key, Record: pro}
		results = append(results, queryResult)

	}
	proAsBytes, err := json.Marshal(results)
	return shim.Success(proAsBytes)
}

//query information of production or person 查询个人信息
func (t *IntelProChainCode) queryuserinfo(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {

	if len(args) != 1 {
		return shim.Error("Please input the correct number of production you want to query ")
	}

	proinfoBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(proinfoBytes)
}

//查询个人信息
func (t *IntelProChainCode) querycontinfo(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {

	if len(args) != 1 {
		return shim.Error("Please input the correct number of content you want to query ")
	}

	proinfoBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(proinfoBytes)
}

//usage of production
func (t *IntelProChainCode) subscribe(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {
	if len(args) != 3 {
		return shim.Error("Please input the correct format(pro number,per number,purpose)")
	}
	var sum float64
	var appinfo AppInfo
	var perinfo Person

	proinfoBytes, _ := APIstub.GetState(args[0])
	_ = json.Unmarshal(proinfoBytes, &appinfo)

	perinfoBytes, _ := APIstub.GetState(appinfo.Applicantid)
	_ = json.Unmarshal(perinfoBytes, &perinfo)
	//Distribute reward in proportion to the contribution to the owner 将酬劳按照贡献比例分配给产品拥有者
	sum, _ = strconv.ParseFloat(appinfo.Usingprice, 64)

	if args[2] == "0" { // 0 means non-commerical purpose
		sum = sum * 0.1
	}

	a, _ := strconv.ParseFloat(perinfo.Account, 64) //转换为float
	percent, _ := strconv.ParseFloat(appinfo.Contribution, 64)
	m := strconv.FormatFloat(a+sum*percent*0.01, 'f', 2, 64)
	perinfo.Account = m

	perinfoRes, _ := json.Marshal(perinfo)
	APIstub.PutState(appinfo.Applicantid, perinfoRes)
	parnumber := len(appinfo.PartnerInfo)
	//Distribute reward in proportion to the contribution to the owner's partner
	for i := 0; i < parnumber; i++ {

		perinfoBytes, _ = APIstub.GetState(appinfo.PartnerInfo[i].Id)
		_ = json.Unmarshal(perinfoBytes, &perinfo)

		b, _ := strconv.ParseFloat(perinfo.Account, 64)
		p, _ := strconv.ParseFloat(appinfo.PartnerInfo[i].Con, 64)
		n := strconv.FormatFloat(float64(b)+sum*p*float64(0.01), 'f', 2, 64)
		perinfo.Account = n
		perinfoRes, _ = json.Marshal(perinfo)
		APIstub.PutState(appinfo.PartnerInfo[i].Id, perinfoRes)

	}
	//Deduct money from user account
	perinfoBytes, _ = APIstub.GetState(args[1])
	_ = json.Unmarshal(perinfoBytes, &perinfo)
	n, _ := strconv.ParseFloat(perinfo.Account, 64)
	if n-sum < 0 {
		return shim.Error("no enough money to use this production")
	}
	perinfo.Account = strconv.FormatFloat(n-sum, 'f', 2, 64)
	perinfoRes, _ = json.Marshal(perinfo)
	APIstub.PutState(args[1], perinfoRes)

	return shim.Success([]byte("sucessful subscribe the production "))
}

//change the owner of the production
func (t *IntelProChainCode) buy(APIstub shim.ChaincodeStubInterface, args []string) ip.Response {
	if (len(args)-7)%3 != 0 && len(args) != 6 || len(args) == 7 {
		return shim.Error("Please input the correct format (pro number,your information)")
	}
	var appinfo AppInfo
	var perinfo Person
	var parinfo PartInfo
	var sum float64
	parnumber := len(args)
	proinfoBytes, _ := APIstub.GetState(args[0])
	_ = json.Unmarshal(proinfoBytes, &appinfo) //appinfo means preowner information
	if appinfo.Purprice == "" {
		return shim.Error("You can not buy this production")
	}
	//increase the account of preowner
	sum, _ = strconv.ParseFloat(appinfo.Purprice, 64)
	perinfoBytes, _ := APIstub.GetState(appinfo.Applicantid)
	_ = json.Unmarshal(perinfoBytes, &perinfo)
	a, _ := strconv.ParseFloat(perinfo.Account, 64)
	percent, _ := strconv.ParseFloat(appinfo.Contribution, 64)
	perinfo.Account = strconv.FormatFloat(a+sum*0.01*percent, 'f', 2, 64)
	perinfoRes, _ := json.Marshal(perinfo)
	APIstub.PutState(appinfo.Applicantid, perinfoRes)
	parnumber = len(appinfo.PartnerInfo)
	for i := 0; i < parnumber; i++ {

		perinfoBytes, _ = APIstub.GetState(appinfo.PartnerInfo[i].Id)
		_ = json.Unmarshal(perinfoBytes, &perinfo)

		b, _ := strconv.ParseFloat(perinfo.Account, 64)
		p, _ := strconv.ParseFloat(appinfo.PartnerInfo[i].Con, 64)
		n := strconv.FormatFloat(b+sum*p*0.01, 'f', 2, 64)
		perinfo.Account = n
		perinfoRes, _ = json.Marshal(perinfo)
		APIstub.PutState(appinfo.PartnerInfo[i].Id, perinfoRes)

	}

	//change the owner of the production        appinfo means　现在的owner (apply)information
	appinfo.Applicantname = args[1]
	appinfo.Applicantid = args[2]
	appinfo.Usingprice = args[3]
	appinfo.Purprice = args[4]
	appinfo.Hashvalue = args[9]
	appinfo.PartnerInfo = []PartInfo{}
	perinfoBytes, _ = APIstub.GetState(appinfo.Applicantid)
	_ = json.Unmarshal(perinfoBytes, &perinfo)
	a, _ = strconv.ParseFloat(perinfo.Account, 64)

	perinfoBytes, _ = APIstub.GetState(appinfo.Applicantid)
	_ = json.Unmarshal(perinfoBytes, &perinfo)

	//buy separately(no need to input applicationPercentage)
	if parnumber == 5 {
		appinfo.Contribution = strconv.Itoa(100)
		appinfoRes, _ := json.Marshal(appinfo)
		APIstub.PutState(args[0], appinfoRes)

		if a-sum < 0 {
			return shim.Error("no enough money to buy this production")
		}
		perinfo.Account = strconv.FormatFloat(a-sum, 'f', 2, 64)
		perinfoRes, _ := json.Marshal(perinfo)
		APIstub.PutState(appinfo.Applicantid, perinfoRes)
		return shim.Success([]byte("sucessful buy the production "))
	}
	//joint purchase,distribute the share
	appinfo.Contribution = args[5]
	percent, _ = strconv.ParseFloat(appinfo.Contribution, 64)
	if a-sum*percent*0.01 < 0 {
		return shim.Error("no enough money to buy this production")
	}

	perinfo.Account = strconv.FormatFloat(a-sum*percent*0.01, 'f', 2, 64)
	perinfoRes, _ = json.Marshal(perinfo)
	APIstub.PutState(appinfo.Applicantid, perinfoRes)

	for i := 6; i < len(args)-1; {
		parinfo.Name = args[i]
		parinfo.Id = args[i+1]
		parinfo.Con = args[i+2]
		perinfoBytes, _ = APIstub.GetState(parinfo.Id)
		_ = json.Unmarshal(perinfoBytes, &perinfo)
		b, _ := strconv.ParseFloat(perinfo.Account, 64)
		p, _ := strconv.ParseFloat(parinfo.Con, 64)
		if float64(b)-sum*p*float64(0.01) < 0 {
			return shim.Error("no enough money to buy this production")
		}
		appinfo.PartnerInfo = append(appinfo.PartnerInfo, parinfo)

		perinfo.Account = strconv.FormatFloat(float64(b)-sum*p*float64(0.01), 'f', 2, 64)
		perinfoRes, _ = json.Marshal(perinfo)
		APIstub.PutState(parinfo.Id, perinfoRes)

		i = i + 3

	}

	appinfores, err := json.Marshal(appinfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	APIstub.PutState(args[0], appinfores)

	return shim.Success([]byte("sucessful buy the production "))

}
func main() {
	err := shim.Start(new(IntelProChainCode))
	if err != nil {
		fmt.Printf("Error starting IntellectualProperty chaincode: %s ", err)
	}
}

