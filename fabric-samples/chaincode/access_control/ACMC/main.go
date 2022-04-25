/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

         http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
    "fmt"
    "encoding/json"
    "strings"
    "crypto/sha1"
    //"crypto/x509"
    //"encoding/pem"
    //"log"
    //"os"
    "time"
    //"os/exec"
    //"bytes"
    "strconv"
    //"math/rand"
    "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("ex02 Init")
    var err error
    Resourcex := map[string]string{"CID":"ffeeab2343feef","SubkeyInfo":"RequiredNum:3,addr1:hash(Subkey1),addr2:hash(Subkey2),addr3:hash(Subkey3)","Owner":"user1","Dept":"FD"}
    B := "cid~owner~txID"
    str, err := json.Marshal(Resourcex)
    if err != nil {
        fmt.Println(err)
    }
    err = stub.PutState(B, []byte(str))
    if err != nil {
        return shim.Error(err.Error())
    }
    //PolicyDatabase := map[string ]map[string] map[string] map[string] string{}
    SA := map[string]string{"Role":"*","Level":"5"}
    OA := map[string]string{"ResourceID":"resource1"}
    AA := map[string]string{"OP":"Download"}
    JudgmentType := map[string]string{"RCA":"PermitOverrides"}
    Target := map[string] map[string] string{}
    Target["SA"]=SA
    Target["OA"]=OA
    Target["AA"]=AA
    Target["JudgmentType"]=JudgmentType
    Action1 := map[string]string{"Efect":"Permit","OP":"AND"}
    Functions1 := map[string]string{"FuncNum":"2","Function1":"isSameDepartment","Args1":"SA,OA","Function2":"more_than","Args2":"SA.Level,2"}
    Rule1 := map[string]map[string]string{}
    Rule1["Action"]=Action1
    Rule1["Functions"]=Functions1
    Action2 := map[string]string{"Efect":"Permit","OP":"OR"}
    Functions2 := map[string]string{"FuncNum":"2","Function1":"isOwner","Args1":"SA,OA.Owner","Function2":"isManager","Args2":"SA.Role"}
    Rule2 := map[string]map[string]string{}
    Rule2["Action"]=Action2
    Rule2["Functions"]=Functions2
    Action3 := map[string]string{"Efect":"Deny","OP":"OR"}
    Functions3 := map[string]string{"FuncNum":"1","Function1":"string_equal","Args1":"SA.Role,guest"}
    Rule3 := map[string]map[string]string{}
    Rule3["Action"]=Action3
    Rule3["Functions"]=Functions3

    Resource1 := map[string]map[string]map[string]string{}
    Resource1["Target"]=Target
    Resource1["Rule1"]=Rule1
    Resource1["Rule2"]=Rule2
    Resource1["Rule3"]=Rule3
    //PolicyDatabase["resource1"]=Resource1
    //A:= "policy_database"
    A := "policyID~txID"
    str_policy, err := json.Marshal(Resource1)
    if err != nil {
        fmt.Println(err)
    }
    err = stub.PutState(A, []byte(str_policy))
    if err != nil {
        return shim.Error(err.Error())
    }
    /*record_map := map[string]string{"Requester":"","ResourceID":"","GrantTime":"","Action":"","DecisionResult":"Permit","EffevtiveDuration":"2h"}
    record_str,err := json.Marshal(record_map)
    if err != nil {
        fmt.Println(err)
    }
    sha1Inst := sha1.New()
    sha1Inst.Write(record_str)
    record_hash := sha1Inst.Sum([]byte(""))
    record_index := string(fmt.Sprintf("%x",record_hash))
    err = stub.PutState(record_index, []byte(record_str))
    if err != nil {
        return shim.Error(err.Error())
    }*/
    fmt.Println("ex02 Init2")
    return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("ex02 Invoke")
    function, args := stub.GetFunctionAndParameters()
    if function == "getToken" {
        return t.getToken(stub, args)
    }else if function == "queryObject" {
        return t.queryObject(stub,args)
    }else if function == "queryPolicy" {
        return t.queryPolicy(stub,args)
    }else if function == "addResource" {
        return t.addResource(stub,args)
    }else if function == "addPolicy" {
        return t.addPolicy(stub,args)
    }else if function == "Read" {
        return t.Read(stub,args)
    }else if function == "Write" {
        return t.Write(stub,args)
    }else if function == "HighThroughputWrite" {
        return t.HighThroughputWrite(stub,args)
    }else if function == "HighThroughputRead" {
        return t.HighThroughputRead(stub,args)
    }
    return shim.Error("Invalid invoke function name. Expecting \"grantAccess\"")
}

func (t *SimpleChaincode) HighThroughputRead(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	A = args[0]
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("varName~txID", []string{A})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", A))
    }
    responseRange, _ := deltaResultsIterator.Next()
    Avalbytes:=responseRange.Value
	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}
func (t *SimpleChaincode) HighThroughputWrite(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string    // Entities
	var Aval  string // Asset holdings
	var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	A = args[0]
    Aval = args[1]
    txid := stub.GetTxID()
    compositeIndexName := "varName~txID"
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{A,txid})
    if compositeErr != nil {
        return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", A, compositeErr.Error()))
    }
	err=stub.PutState(compositeKey, []byte(Aval))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string    // Entities
	var Aval  string // Asset holdings
	var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	A = args[0]
    Aval = args[1]
	// Write the state back to the ledger
	err = stub.PutState(A, []byte(Aval))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func (t *SimpleChaincode) Read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	A = args[0]
	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func (t *SimpleChaincode) addPolicy(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	var err error
    var err_flag int
    var subject_attr map[string]string
    subject_attr = make(map[string]string)
    subject_attr,err_flag = getSubject(stub) //通过调用pip链码获取属性
    if err_flag == 1{
        return shim.Success([]byte("getSubject false!"))
    }
    SA := map[string]string{"Role":args[0],"Level":args[1]}
    OA := map[string]string{"ResourceID":args[2]}
    var object_attr map[string]string
    object_attr = make(map[string]string)
    tmp_args := []string{args[2]}
    object_attr,err_flag = getObject(stub,tmp_args)
    if err_flag == 1{
        return shim.Success([]byte("Get Object False!"))
    }
    if strings.Compare(subject_attr["Name"],object_attr["Owner"])!=0{
		return shim.Success([]byte("Permit Deny!"))
    }
    AA := map[string]string{"OP":args[3]}
    JudgmentType := map[string]string{"RCA":args[4]}
    Target := map[string] map[string] string{}
    Target["SA"]=SA
    Target["OA"]=OA
    Target["AA"]=AA
    Target["JudgmentType"]=JudgmentType
    Action1 := map[string]string{"Efect":args[5],"OP":args[6]}
    Functions1 := map[string]string{"FuncNum":args[7],"Function1":args[8],"Args1":args[9],"Function2":args[10],"Args2":args[11]}
    Rule1 := map[string]map[string]string{}
    Rule1["Action"]=Action1
    Rule1["Functions"]=Functions1
    Action2 := map[string]string{"Efect":args[12],"OP":args[13]}
    Functions2 := map[string]string{"FuncNum":args[14],"Function1":args[15],"Args1":args[16],"Function2":args[17],"Args2":args[18]}
    Rule2 := map[string]map[string]string{}
    Rule2["Action"]=Action2
    Rule2["Functions"]=Functions2
    Action3 := map[string]string{"Efect":args[19],"OP":args[20]}
    Functions3 := map[string]string{"FuncNum":args[21],"Function1":args[22],"Args1":args[23]}
    Rule3 := map[string]map[string]string{}
    Rule3["Action"]=Action3
    Rule3["Functions"]=Functions3

    Resource1 := map[string]map[string]map[string]string{}
    Resource1["Target"]=Target
    Resource1["Rule1"]=Rule1
    Resource1["Rule2"]=Rule2
    Resource1["Rule3"]=Rule3
    str1,err := json.Marshal(Resource1)
    if err != nil {
        fmt.Println(err)
    }
    txid := stub.GetTxID()
    compositeIndexName := "policyID~txID"
    policyID:="policy"+args[2]
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{policyID,txid})
    if compositeErr != nil {
        return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", policyID, compositeErr.Error()))
    }
	err=stub.PutState(compositeKey, []byte(str1))
	if err != nil {
		return shim.Error(err.Error())
	}
        return shim.Success([]byte(fmt.Sprintf("Add %s Successfully!",compositeKey)))
}

func (t *SimpleChaincode) addResource(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var err_flag int
    var subject_attr map[string]string
    subject_attr = make(map[string]string)
    subject_attr,err_flag = getSubject(stub) //通过调用pip链码获取属性
    if err_flag == 1{
        return shim.Error("getSubject false!")
    }
	//var err error
	var cid,dept,owner,subkeyinfo string  // Entities
    owner = subject_attr["Name"]
    cid = args[0]
    dept = args[1]
    subkeyinfo = args[2]
    var policy_args = []string{cid}
    _,err_flag = getPolicy(stub,policy_args)
    if err_flag==1{
		//return shim.Error("Failed to get policy!")
        fmt.Println("Policy is not found!!")
    }
    resource := map[string]string{"CID":cid,"Owner":owner,"Dept":dept,"SubkeyInfo":subkeyinfo}
    str1,err := json.Marshal(resource)
    if err != nil {
        fmt.Println(err)
    }
    txid := stub.GetTxID()
    compositeIndexName := "cid~owner~txID"
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{cid,owner,txid})
    if compositeErr != nil {
        return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", cid, compositeErr.Error()))
    }
	err=stub.PutState(compositeKey, []byte(str1))
	if err != nil {
		return shim.Error(err.Error())
	}
    resourceID_str := fmt.Sprintf("resource:%s added!",cid)
    return shim.Success([]byte(resourceID_str))
}
func (t *SimpleChaincode) queryPolicy(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var resource map[string]map[string]map[string]string
    resource = make(map[string]map[string]map[string]string)
	var err error
    var policyID string
    policyID = "policy"+args[0]
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("policyID~txID", []string{policyID})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    //defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", policyID))
    }
    responseRange, _ := deltaResultsIterator.Next()
    err = json.Unmarshal(responseRange.Value, &resource)
    if err != nil {
        fmt.Println(err)
    }
    str1,err := json.Marshal(resource)
    if err != nil {
        fmt.Println(err)
    }
	return shim.Success(str1)
}
func (t *SimpleChaincode) queryObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var resourceID string
    var err error
    var  resource_attr  map[string]string
    resource_attr = make(map[string]string)
    resourceID = args[0]
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("cid~owner~txID", []string{resourceID})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", resourceID))
    }
    responseRange, _ := deltaResultsIterator.Next()
    err = json.Unmarshal(responseRange.Value, &resource_attr)
    if err != nil {
        fmt.Println(err)
    }
    str1,err := json.Marshal(resource_attr)
    if err != nil {
        fmt.Println(err)
    }
	return shim.Success(str1)
}
// PIP functions
func  getEnvirment(stub shim.ChaincodeStubInterface) string{
    timestamp, err  := stub.GetTxTimestamp()
    if err != nil{
        return "false"
    }
    tm:= time.Unix(timestamp.Seconds, 0)
    ts := tm.Format("2006-01-02 03:04:05 PM")
    return string(ts)
}

func  getObject(stub shim.ChaincodeStubInterface, args []string) (map[string]string,int) {
    var resourceID string
    flag := 1
    var err error
    var resource_attr map[string]string
    resource_attr = make(map[string]string) //LNB
    resourceID = args[0] //CID
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("cid~owner~txID", []string{resourceID})
    if deltaErr != nil {
        return resource_attr,flag
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        //return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
        return resource_attr,flag
    }
    responseRange, nextErr := deltaResultsIterator.Next()
    if nextErr != nil {
        return resource_attr,flag
    }
    err = json.Unmarshal(responseRange.Value, &resource_attr)
    if err != nil {
        fmt.Println(err)
    }
    flag=0
    return resource_attr,flag
}


func  getSubject(stub shim.ChaincodeStubInterface) (map[string]string,int){
    flag := 1
    subject_attr := map[string]string{}
    val,ok,err := cid.GetAttributeValue(stub, "Role")
    if err != nil {
        return subject_attr,flag
    }
    if !ok {
        return subject_attr,flag
    }
    //subject_attr["Role"]=string(val)
    val1,ok,err := cid.GetAttributeValue(stub, "Department")
    if err != nil {
        return subject_attr,flag
    }
    if !ok {
        return subject_attr,flag
    }
    //subject_attr["Department"]=string(val1)
    val2, ok, err := cid.GetAttributeValue(stub, "Level")
    if err != nil {
        return subject_attr,flag
    }
    if !ok {
        return subject_attr,flag
    }
    val3,ok,err := cid.GetAttributeValue(stub, "Name")
    if err != nil {
        return subject_attr,flag
    }
    if !ok {
        return subject_attr,flag
    }
    subject_attr["Role"]=val
    subject_attr["Department"]=val1
    subject_attr["Level"]=val2
    subject_attr["Name"]=val3
    flag = 0
    return subject_attr,flag
}

func isSameDepartment(sa_map map[string]string,oa_map map[string]string) bool{
    if strings.Compare(sa_map["Department"],oa_map["Dept"])==0{
        return true
    }
    return false
}

func more_than(subject_level string, level int) bool {
    i, err := strconv.Atoi(subject_level)
    if err != nil {
        panic(err)
    }
    if i > level{
        return true
    }
    return false
}

func isManager(role string) bool{
    if strings.Compare(role,"manager")==0{
        return true
    }
    return false
}
func isOwner(user string,owner string) bool{
    if strings.Compare(user,owner)==0 {
        return true 
    }
    return false
}
func isMeetedTarget() bool{
    return true
}
func ANDfunction(val1 bool,val2 bool) bool{
    if val1==true && val2 == true{
        return true
    }
    return false
}

func ORfunction(val1 bool,val2 bool) bool{
    if val1==false && val2 == false{
        return false
    }
    return true
}
func permitOverrides(rule1 string ,rule2 string ,rule3 string) string{
    permit:="Permit"
    deny:="Deny"
    
    if rule1==permit || rule2==permit || rule3==permit{
        return permit
    }
    return deny
}
func denyOverrides(rule1 string ,rule2 string ,rule3 string) string{
    permit:="Permit"
    deny:="Deny"
    if rule1==deny || rule2==deny || rule3==deny{
        return deny
    }
    return permit
}

func string_equal(str1 string,str2 string) bool{
    if strings.Compare(str1,str2)==0{
        return true
    }
    return false
}

func  getDecision(stub shim.ChaincodeStubInterface,args []string) string{
    var err_flag int
    var err error
    var attrs string
    var attrs_map map[string] map[string] string
    attrs_map = make(map[string] map[string] string)
    var sa_map map [string] string
    sa_map = make(map[string] string)
    var oa_map map [string] string
    oa_map = make(map[string] string)
    var ea_map map [string] string
    ea_map = make(map[string] string)
    var aa_map map [string] string
    aa_map = make(map[string] string)
    attrs = args[0]
    err = json.Unmarshal([]byte(attrs), &attrs_map)
    if err != nil {
        fmt.Println(err)
    }
    sa_map=attrs_map["SA"]
    oa_map=attrs_map["OA"]
    ea_map=attrs_map["EA"]
    fmt.Println(ea_map)
    aa_map=attrs_map["AA"]
    fmt.Println(aa_map)
    resourceID := oa_map["CID"]
    var policy map[string] map[string] map[string] string
    policy = make(map[string] map[string] map[string] string)
    var policy_args = []string{resourceID}
    policy,err_flag = getPolicy(stub,policy_args)
    if err_flag==1{
        return "false"
    }
    fmt.Println(policy)
    /*err = json.Unmarshal([]byte(policy_bytes.Payload), &policy)
    if err != nil {
        fmt.Println(err)
    }*/
    if isMeetedTarget() == false{
        return "false"
    }
    var rule1_val string
    var rule2_val string
    var rule3_val string
    if ANDfunction(isSameDepartment(sa_map,oa_map),more_than(sa_map["Level"],0)) == true{
        rule1_val = "Permit"
    }else{
        rule1_val = "Dissatisfy"
    }

    if ORfunction(isOwner(sa_map["Name"],oa_map["Owner"]),isManager(sa_map["Role"])) == true{
        rule2_val = "Permit"
    }else{
        rule2_val = "Dissatisfy"
    }
    if string_equal(sa_map["Role"],"guest") == true{
        rule3_val = "Deny"
    }else{
        rule3_val = "Dissatisfy"
    }
    result := permitOverrides(rule1_val,rule2_val,rule3_val)
    /*policy_str,err := json.Marshal(policy)
    if err != nil {
        fmt.Println(err)
    }
    return shim.Success(policy_str)*/
    return result
}
//PAP Functions
func  getPolicy(stub shim.ChaincodeStubInterface,args []string) (map[string]map[string]map[string]string,int){
    flag := 1
    var resource map[string]map[string]map[string]string
    resource = make(map[string]map[string]map[string]string)
	var err error
    var policyID string
    policyID = "policy"+args[0]
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("policyID~txID", []string{policyID})
    if deltaErr != nil {
		//return shim.Error("Failed to get state by partial compositeKey!")
        return resource,flag
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return resource,flag
        //return shim.Error(fmt.Sprintf("No variable by the name %s exists", policyID))
    }
    responseRange, _ := deltaResultsIterator.Next()
    err = json.Unmarshal(responseRange.Value, &resource)
    if err != nil {
        fmt.Println(err)
    }
    flag = 0
    return resource,flag 
}

//PEP Functions
func (t *SimpleChaincode) getToken(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var err_flag int
    var err error
    var resourceID,action string
    var all_attr map[string]map[string]string
    all_attr = make(map[string] map[string] string)
    resourceID = args[0]
    fmt.Println(resourceID)
    action = args[1]
    var subject_attr map[string]string
    subject_attr = make(map[string]string)
    subject_attr,err_flag = getSubject(stub) //通过调用pip链码获取属性
    if err_flag == 1{
        return shim.Error("getSubject false!")
    }
    var object_attr map[string]string
    object_attr = make(map[string]string)
    object_attr,err_flag = getObject(stub,args)
    if err_flag == 1{
        return shim.Error("getObject false!")
    }
    var env_attr map[string]string
    env_attr = make(map[string]string)
    env_attr["Time"] = getEnvirment(stub)
    if env_attr["Time"]=="false"{
        return shim.Error("getEnvirment false!")
    }
    var action_attr map[string] string
    action_attr = make(map[string] string)
    action_attr["OP"]=action
    all_attr["SA"]=subject_attr
    all_attr["OA"]=object_attr
    all_attr["AA"]=action_attr
    all_attr["EA"]=env_attr  //all attributevalue 组建完成
    attrs,err := json.Marshal(all_attr)
    if err != nil {
        return shim.Error("json.Marshal(all_attr) false!")
    }
// invoke pdp 
    var getDecision_args = []string{string(attrs)}
    decision_result := getDecision(stub,getDecision_args)
    if decision_result == "false"{
        return shim.Error("getDecision false")
    }
    record_map := map[string]string{"Requester":subject_attr["Name"],"CID":object_attr["CID"],"GrantTime":env_attr["Time"],"Action":action,"DecisionResult":decision_result}
    result_map := map[string]string{"DecisionResult":decision_result}
    if decision_result == "Permit"{
        record_map["EffevtiveDuration"]="2"
        //record_map["SubkeyInfo"]=object_attr["SubkeyInfo"]
        //record_map["CID"]=object_attr["CID"]

    }else{
        record_map["EffevtiveDuration"]="0"
        //record_map["SubkeyInfo"]="Null"
        //record_map["CID"]="Null"
    }
    str_record,err := json.Marshal(record_map)
    if err != nil {
        return shim.Error("json.Marshal(record_map) false!")
    }
    sha1Inst := sha1.New()
    sha1Inst.Write(str_record)
    record_hash := sha1Inst.Sum([]byte(""))
    tokenID := string(fmt.Sprintf("%x",record_hash))

    ar_args := [][]byte{[]byte("putGrantLog"),str_record,[]byte(tokenID)}
    response_PutGrand := stub.InvokeChaincode("ALMC",ar_args,"mychannel") //通过调用pdp链码获取决策结果
    if response_PutGrand.Status !=shim.OK{
        return shim.Error("failed to invoke putGrantRecord!")
    }
    if decision_result == "Permit"{
        result_map["Requester"]=subject_attr["Name"]
        result_map["SubkeyInfo"]=object_attr["SubkeyInfo"]
        result_map["TokenID"]=tokenID
        result_map["DecisionResult"]="Permit"
        result_map["CID"]=object_attr["CID"]
        str_result,err := json.Marshal(result_map)
        if err != nil {
            return shim.Error("json.Marshal(result_map) false!")
        }
        return shim.Success(str_result)
    }else{
        //return shim.Error("Authorization Failure!")
        return shim.Success([]byte("Authorization Failure!"))
    }
}

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

