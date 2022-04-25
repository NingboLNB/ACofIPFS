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
    "encoding/hex"
    //"crypto/sha1"
	"crypto/ecdsa"
	//"crypto/rand"
    "math/big"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	//"io/ioutil"
    //"strings"
    //"crypto/x509"
    //"encoding/pem"
    //"log"
    //"os"
    "time"
    //"os/exec"
    //"bytes"
    "strconv"
    //"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
    //"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	//_, args := stub.GetFunctionAndParameters()
	var err error
    grantlog := map[string]string{"Requester":"","Action":"","DecisionResult":"","EffevtiveDuration":"0","ResourceID":"","GrantTime":""}
    accesslog := map[string]map[string]string{"AccessTime":{"tokenID":"","Info":""}}
    grantlog_str,err := json.Marshal(grantlog)
    if err != nil {
        fmt.Println(err)
    }
	err = stub.PutState("grantlog~tokenID~txID", grantlog_str)
	if err != nil {
		return shim.Error(err.Error())
	}
    accesslog_str,err := json.Marshal(accesslog)
    if err != nil {
        fmt.Println(err)
    }
	err = stub.PutState("accesslog~requester~txID", accesslog_str)
	if err != nil {
		return shim.Error(err.Error())
	}
    /*sha1Inst := sha1.New()
    sha1Inst.Write(record_str)
    record_hash := sha1Inst.Sum([]byte(""))
    record_index := string(fmt.Sprintf("%x",record_hash))
	err = stub.PutState(record_index, []byte(record_str))
	if err != nil {
		return shim.Error(err.Error())
	}*/
    pk_database := map[string]string{"userx":"certBytes"}
    pk_str,err := json.Marshal(pk_database)
    if err != nil {
        fmt.Println(err)
    }
	err = stub.PutState("pk_database", pk_str)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func isvalid(str_record []byte) bool {
    grantlog := map[string]string{}
    err := json.Unmarshal(str_record, &grantlog)
    if err != nil {
        fmt.Println(err) 
    }
    GrantTime:=grantlog["GrantTime"]
    fmt.Println(GrantTime)
    CurrentTime:=time.Now().Format("2006-01-02 03:04:05")
    fmt.Println(CurrentTime)
    Duration:=grantlog["EffevtiveDuration"]
    flag,_ :=strconv.Atoi(Duration)
    if flag > 0{
        return true
    }
    return false
}
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "putGrantLog" {
		return t.putGrantLog(stub, args)
	} else if function == "tokenValidation" {
		return t.tokenValidation(stub,args)
	} else if function == "putPublicKey" {
		return t.putPublicKey(stub,args)
	} else if function == "queryAccessLog" {
		return t.queryAccessLog(stub,args)
	} else if function == "queryGrantLog" {
		return t.queryGrantLog(stub,args)
    }
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) tokenValidation(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var err error
    grantlog:=map[string]string{}
    requster:=args[0]
    tokenID:=args[1]
    strSigR:=args[2]
    strSigS:=args[3]
    rint := big.Int{}
	sint := big.Int{}
	rByte, _ := hex.DecodeString(strSigR)
	sByte, _ := hex.DecodeString(strSigS)
	rint.SetBytes(rByte)
	sint.SetBytes(sByte)
	//fmt.Println("------", rint.SetBytes(rByte))
	//fmt.Println("------", sint.SetBytes(sByte))
	hash := sha256.Sum256([]byte(tokenID))

    /*A := "pk_database"
    map2 := map[string]string{}
    str,err := stub.GetState(A)
    if err != nil{
        return shim.Error("get pubkey failed!")
    }
    err = json.Unmarshal(str, &map2)
    if err != nil {
        return shim.Error("unmarshal Failed!")
    }
    certBytes := []byte(map2[requster])*/
    
    pk_username:="pk_"+requster
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("pk_username~txID", []string{pk_username})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", pk_username))
    }
    responseRange, _ := deltaResultsIterator.Next()
    certBytes := responseRange.Value
	blkCert, _ := pem.Decode(certBytes)
	//fmt.Println("cert.pem type:", blkCert.Type)
	cert, _ := x509.ParseCertificate(blkCert.Bytes)
	pubkey := cert.PublicKey.(*ecdsa.PublicKey)
    ok := ecdsa.Verify(pubkey, hash[:], &rint, &sint)
    if ok==true{
        GLtokenID:="GL"+tokenID
        deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("GLtokenID~txID", []string{GLtokenID})
        if deltaErr != nil {
		    return shim.Error("Failed to get state by partial compositeKey!")
        }
        defer deltaResultsIterator.Close()
        if !deltaResultsIterator.HasNext() {
            info:="Token Not Found"
            tmp_args := []string{requster,tokenID,info}
            putAccessLog(stub,tmp_args)
            return shim.Success([]byte(info))
        }
        responseRange, _ := deltaResultsIterator.Next()
        err = json.Unmarshal(responseRange.Value, &grantlog)
        if err != nil {
            fmt.Println(err)
        }
        str_record:=responseRange.Value
        if isvalid(str_record){
            info:="Success!"
            tmp_args := []string{requster,tokenID,info}
            putAccessLog(stub,tmp_args)
            return shim.Success([]byte(info))
        }else{
            info:="Failure,Token Expired!"
            tmp_args := []string{requster,tokenID,info}
            putAccessLog(stub,tmp_args)
            return shim.Success([]byte("Failure,Token Expired!"))
        }
    }else{
        info:="Invalid Signature!"
        tmp_args := []string{requster,tokenID,info}
        putAccessLog(stub,tmp_args)
         // return shim.Success(certBytes)
		return shim.Error(info)
        //return shim.Success([]byte(info))
    }
}
func (t *SimpleChaincode) putPublicKey(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var err error
    /*var username string
    var pk string
    username=args[0]
    pk=args[1]
    A := "pk_database"
    map2 := map[string]string{}
    str,_ := stub.GetState(A)
    err = json.Unmarshal(str, &map2)
    if err != nil {
        fmt.Println(err)
    }
    map2[username]=pk
    str1,err := json.Marshal(map2)
    fmt.Println(string(str1))
    if err != nil {
        fmt.Println(err)
    }
	err=stub.PutState(A,str1)
	if err != nil {
		return shim.Error(err.Error())
	}*/
    username:=args[0]
    cert:=args[1]
    pk_username:="pk_"+username
    txid := stub.GetTxID()
    compositeIndexName := "pk_username~txID"
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{pk_username,txid})
    if compositeErr != nil {
        return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", username, compositeErr.Error()))
    }
	err=stub.PutState(compositeKey, []byte(cert))
	if err != nil {
		return shim.Error(err.Error())
	}

    return shim.Success(nil)
}
func  getEnvirment(stub shim.ChaincodeStubInterface) string{
    timestamp, err  := stub.GetTxTimestamp()
    if err != nil{
        return "false"
    }
    tm:= time.Unix(timestamp.Seconds, 0)
    ts := tm.Format("2006-01-02 03:04:05 PM")
    return string(ts)
}
func (t *SimpleChaincode) queryAccessLog(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    username:=args[0]
    AL_requester:="AL_"+username
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("AL_requester~txID", []string{AL_requester})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", AL_requester))
    }
    responseRange, _ := deltaResultsIterator.Next()
    return shim.Success(responseRange.Value)
}
func (t *SimpleChaincode) queryGrantLog(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    tokenID:=args[0]
    GLtokenID:="GL"+tokenID
    deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("GLtokenID~txID", []string{GLtokenID})
    if deltaErr != nil {
		return shim.Error("Failed to get state by partial compositeKey!")
    }
    defer deltaResultsIterator.Close()
    if !deltaResultsIterator.HasNext() {
        return shim.Error(fmt.Sprintf("No variable by the name %s exists", GLtokenID))
    }
    responseRange, _ := deltaResultsIterator.Next()
    return shim.Success(responseRange.Value)
}
func putAccessLog(stub shim.ChaincodeStubInterface,args []string) int{
    var err error
    AccessTime:=getEnvirment(stub)
    username:=args[0]
    tokenID:=args[1]
    info:=args[2]
    user_map := map[string]map[string]string{}
    map2 := map[string]string{"tokenID":tokenID,"Info":info}
    user_map[AccessTime] = map2
    accesslog_str,err := json.Marshal(user_map)
    if err != nil {
        fmt.Println(err)
    }
    txid := stub.GetTxID()
    compositeIndexName := "AL_requester~txID"
    AL_requester:="AL_"+ username
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{AL_requester,txid})
    if compositeErr != nil {
        return -1
    }
	err=stub.PutState(compositeKey, []byte(accesslog_str))
	if err != nil {
		return -1
	}
	return 0
}
func (t *SimpleChaincode) putGrantLog(stub shim.ChaincodeStubInterface,args []string) pb.Response{
    var err error
    var record_str string
    record_str = args[0]
    record_index := args[1]
    txid := stub.GetTxID()
    compositeIndexName := "GLtokenID~txID"
    GLtokenID:="GL"+record_index
    compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{GLtokenID,txid})
    if compositeErr != nil {
        return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", record_index, compositeErr.Error()))
    }
	err=stub.PutState(compositeKey, []byte(record_str))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

