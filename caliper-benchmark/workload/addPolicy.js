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
            contractFunction: 'addPolicy',
            invokerIdentity: 'user4@org1.example.com',
            contractArguments: ["*","5","QmY4GC1QqLMV6XQrGpwuqtGD3wC3nm4oJcbxT7N","Download","PermitOverrides","Permit","AND","2","isSameDepartment","SA,OA","more_than","SA.Level,2","Permit","OR","2","isOwner","SA,OA.Owner","isManager","SA.Role","Deny","OR","1","string_equal","SA.Role,guest"],
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
    
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;


