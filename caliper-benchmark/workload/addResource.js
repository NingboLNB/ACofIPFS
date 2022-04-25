'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }
    
    
    async submitTransaction() {
        const randomId = Math.floor(Math.random()*this.roundArguments.assets);
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'addResource',
            invokerIdentity: 'user4@org1.example.com',
            contractArguments: ["QmY4GVmEC1QqLMV6XQrGpwuqtGD3wC3nm4oJcbxT7zbWLN","FD","RequiredNum:3,peer0.org1:3595bcb0e9aeb631b07b5067fcd33ea6abfddd83,peer1.org1:1f49fdf0c9e0e7481c2063f7e5af7c25c61f50c8,peer0.org2:c7463526013c7a53b3ef0ad25bffcbe3643d60a7,peer1.org2:ca7594bdbd97b7f30ac140c1fd307c563fac41d0"],
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;


