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
            contractFunction: 'queryObject',
            invokerIdentity: 'user4@org1.example.com',
            contractArguments: ["QmY4GC1QqLMV6XQrGpwuqtGD3wC3nm4oJcbxT7N"],
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
    
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;


