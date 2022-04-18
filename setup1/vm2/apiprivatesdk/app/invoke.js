const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory, DefaultCheckpointers } = require('fabric-network');
const fs = require('fs');
const EventStrategies = require('fabric-network/lib/impl/event/defaulteventhandlerstrategies');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('TracechainNetwork');
const util = require('util')

const helper = require('./helper');
const { blockListener, contractListener, commitListener } = require('./Listeners');
const { resolve } = require('path');
const { fullList } = require('npm');

const submitTransactionWithPrivateData = async (contract, fcn, args1,args2,transientData) => {
	try{
		let transientJSON = JSON.parse(transientData);
	        let key = Object.keys(transientJSON)[0];
		let result;
		const transientDataBuffer = {};
		transientDataBuffer[key] = Buffer.from(JSON.stringify((transientJSON.asset_properties)));
		for(let i=0;i<10;i++){
		try {
		result = await contract.createTransaction(fcn)
			.setTransient(transientDataBuffer)
        		.submit(args1,args2)
		break;
		} catch (error) {
			console.log(`Creation of transaction failed attempting again [Attempt: ${i}] `)
			i++;
		}
		
		}
		if(result !== undefined){
			return result;
		} else {
			console.log(`Invoke Transaction Getting error: ${error}`);
        		const response_payload = {
            			result: null,
           			 error: true
       			 };
        	return response_payload;
		}
	}catch(error){
	console.log(`Invoke Transaction Getting error: ${error}`);
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        };
        return response_payload;
	}
}

const submitTransactionWithUpdatePrivateData = async (contract, fcn, args0, args1, args2, args3, args4, args5, args6, args7, args8, args9, args10, args11, args12,args13, args14, args15, args16, args17, args18, args19, args20, args21, args22,args23,transientData) => {
	try{
		let transientJSON = JSON.parse(transientData);
	        let key = Object.keys(transientJSON)[0];
		let result;
		const transientDataBuffer = {};
		transientDataBuffer[key] = Buffer.from(JSON.stringify((transientJSON.asset_properties)));
		for(let i=0;i<10;i++){
		try {
		result = await contract.createTransaction(fcn)
			.setTransient(transientDataBuffer)
        		.submit(args0, args1, args2, args3, args4, args5, args6, args7, args8, args9, args10, args11, args12,args13, args14, args15, args16, args17, args18, args19, args20, args21, args22,args23)
		break;
		} catch (error) {
			console.log(`Creation of transaction failed attempting again [Attempt: ${i}] `)
			i++;
		}
		
		}
		if(result !== undefined){
			return result;
		} else {
			console.log(`Invoke Transaction Getting error: ${error}`);
        		const response_payload = {
            			result: null,
           			 error: true
       			 };
        	return response_payload;
		}
	}catch(error){
	console.log(`Invoke Transaction Getting error: ${error}`);
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        };
        return response_payload;
	}
}


const submitTransactionWithDeletePrivateData = async (contract, fcn, args0, args1,transientData) => {
	try{
		let transientJSON = JSON.parse(transientData);
	        let key = Object.keys(transientJSON)[0];
		let result;
		const transientDataBuffer = {};
		transientDataBuffer[key] = Buffer.from(JSON.stringify((transientJSON.asset_properties)));
		for(let i=0;i<10;i++){
		try {
		result = await contract.createTransaction(fcn)
			.setTransient(transientDataBuffer)
        		.submit(args0, args1)
		break;
		} catch (error) {
			console.log(`Creation of transaction failed attempting again [Attempt: ${i}] `)
			i++;
		}
		
		}
		if(result !== undefined){
			return result;
		} else {
			console.log(`Invoke Transaction Getting error: ${error}`);
        		const response_payload = {
            			result: null,
           			 error: true
       			 };
        	return response_payload;
		}
	}catch(error){
	console.log(`Invoke Transaction Getting error: ${error}`);
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        };
        return response_payload;
	}
}

const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name, transientData) => {
    try {
        const ccp = await helper.getCCP(org_name);

        const walletPath = await helper.getWalletPath(org_name);
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }


        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: false },
            // eventHandlerOptions: EventStrategies.NONE
	        eventHandlerOptions: {
                commitTimeout: 40000,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            }
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);
        await gateway.disconnect();


        const network = await gateway.getNetwork(channelName);

        const contract = network.getContract(chaincodeName);

        let result;
        	console.log(contract.sendInstantiateProposal)


        switch (fcn) {
            case "InsertAssetRecords":
		 console.log("====== Executing InsertAssetRecords ======")
                //result = await contract.submitTransaction(fcn, args[0]);
                if(parseInt(args[1]) === 0){
                    // Submit transaction with private data 
                    result = await submitTransactionWithPrivateData(contract,fcn,args[0],args[1],transientData)
                } else {
                    //submit tranascation without private data 
                    result = await contract.submitTransaction(fcn,args[0],args[1])
                }
                console.log("#######################",JSON.stringify(result))
                console.log(`Output: InsertAssetRecords= ${result.toString()}`)
                result = {txid: result.toString()}
                break;
            case "InsertOrderRecords":
                console.log("====== Executing InsertOrderRecords ======")
                //result = await contract.submitTransaction(fcn, args[0]);
                if(parseInt(args[1]) === 0){
                    result = await submitTransactionWithPrivateData(contract,fcn,args[0],args[1],transientData)
                } else {
                    result = await contract.submitTransaction(fcn,args[0],args[1])
                }
                console.log(`Output: InsertOrderRecords= ${result.toString()}`)
                result = {txid: result.toString()}
                break;
            case "UpdateAsset":
	        	console.log("====== Executing UpdateAsset ======")
                if(parseInt(args[23]) === 0){
                    result = await submitTransactionWithUpdatePrivateData(contract, fcn, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12],args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22],args[23],transientData);
                } else {
                    result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12],args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22]);
                }
                console.log(`Output: UpdateAsset= ${result.toString()}`)
                result = {txid: result.toString()}
                break;
	       case "DeleteAssetByTransactionId":
                console.log("====== Executing DeleteAssetByTransactionId ======")
                if(parseInt(args[1]) === 0){
                    result = await submitTransactionWithDeletePrivateData(contract, fcn, args[0], args[1], transientData);
                } else {
                    result = await contract.submitTransaction(fcn, args[0], args[1]);
                }
                console.log(`Output: DeleteAssetByTransactionId= ${result.toString()}`)
                result = {txid: result.toString()}
                break;

            default:	
                break;
        }



        await network.addBlockListener(blockListener,{type:"private"})



        const response_payload = {
            result: result,
            error: false,
            errorData: null
        }

        console.log("Invoke Transaction Try Block Response Payload: ", response_payload)
        return response_payload;


    } catch (error) {

        console.log(`Invoke Transaction Getting error: ${error}`)
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        }
        console.log("Invoke Transaction Catch Block Response Payload: ", response_payload)
        return response_payload

    }
}

exports.invokeTransaction = invokeTransaction;
