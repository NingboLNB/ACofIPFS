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
            contractFunction: 'queryGrantLog',
            invokerIdentity: 'user4@org1.example.com',
            contractArguments: ["173ba9a9a5f8dac262af887310a50b4b1efcc9ad"],
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
    
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;


