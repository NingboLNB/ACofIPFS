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
            contractFunction: 'tokenValidation',
            invokerIdentity: 'user4@org1.example.com',
            contractArguments: ["user4","db45f58652d8df5acc3f46ed8ae44d99c85ad751","ab33b719c5934b31c3fe48764b99bf31d809f0f48c23b176688dcb2664e845b0","487e9604414ba00180a4d6200bfe74dade197410fd29b218646322c1bbbc0811"],
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
    
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;


