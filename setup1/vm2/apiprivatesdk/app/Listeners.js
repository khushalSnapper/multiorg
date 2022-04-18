const {BlockListener} = require("fabric-network")
let request = require('request')
const constants = require('../config/constants.json');
let temp =0;
function showTransactionData(transactionData) {
    console.log(JSON.stringify(transactionData))
	const creator = transactionData.actions[0].header.creator;
	console.log(`    - submitted by: ${creator.mspid}-${creator.id_bytes.toString('hex')}`);
	for (const endorsement of transactionData.actions[0].payload.action.endorsements) {
		console.log(`    - endorsed by: ${endorsement.endorser.mspid}-${endorsement.endorser.id_bytes.toString('hex')}`);
	}
	const chaincode = transactionData.actions[0].payload.chaincode_proposal_payload.input.chaincode_spec;
	console.log(`    - chaincode:${chaincode.chaincode_id.name}`);
	console.log(`    - function:${chaincode.input.args[0].toString()}`);
	for (let x = 1; x < chaincode.input.args.length; x++) {
		console.log(`    - arg:${chaincode.input.args[x].toString()}`);
	}
}

// contractListener = async (event) => {
//     console.log("==========================================")
//     console.log(event)
//     // The payload of the chaincode event is the value place there by the
//     // chaincode. Notice it is a byte data and the application will have
//     // to know how to deserialize.
//     // In this case we know that the chaincode will always place the asset
//     // being worked with as the payload for all events produced.
//     const asset = JSON.parse(event.payload.toString());
//     console.log(`<-- Contract Event Received: ${event.eventName} - ${JSON.stringify(asset)}`);
//     // show the information available with the event
//     console.log(`*** Event: ${event.eventName}:${asset.ID}`);
//     // notice how we have access to the transaction information that produced this chaincode event
//     const eventTransaction = event.getTransactionEvent();
//     console.log(`*** transaction: ${eventTransaction.transactionId} status:${eventTransaction.status}`);
//     showTransactionData(eventTransaction.transactionData);
//     // notice how we have access to the full block that contains this transaction
//     const eventBlock = eventTransaction.getBlockEvent();
//     console.log(`*** block: ${eventBlock.blockNumber.toString()}`);
// };


// blockListener = async (event) => {

//     console.log("--------------------------------------------------------------")
//     console.log(`<-- Block Event Received - block number: ${event.blockNumber.toString()}`);

//     const transEvents = event.getTransactionEvents();
//     for (const transEvent of transEvents) {
//         console.log(`*** transaction event: ${transEvent.transactionId}`);
//         // if (transEvent.privateData) {
//         //     for (const namespace of transEvent.privateData.ns_pvt_rwset) {
//         //         console.log(`    - private data: ${namespace.namespace}`);
//         //         for (const collection of namespace.collection_pvt_rwset) {
//         //             console.log(`     - collection: ${collection.collection_name}`);
//         //             if (collection.rwset.reads) {
//         //                 for (const read of collection.rwset.reads) {
//         //                     console.log(`       - read set - ${BLUE}key:${RESET} ${read.key}  ${BLUE}value:${read.value.toString()}`);
//         //                 }
//         //             }
//         //             if (collection.rwset.writes) {
//         //                 for (const write of collection.rwset.writes) {
//         //                     console.log(`      - write set - ${BLUE}key:${RESET}${write.key} ${BLUE}is_delete:${RESET}${write.is_delete} ${BLUE}value:${RESET}${write.value.toString()}`);
//         //                 }
//         //             }
//         //         }
//         //     }
//         // }
//         if (transEvent.transactionData) {
//             showTransactionData(transEvent.transactionData);
//         }
//     }
// };


const blockListener = async (event) => {
    const bdata = event.blockData.data.data;
    
    console.log(bdata);

    const transEvents = event.getTransactionEvents();  
    //console.log(transEvents[0].transactionData.actions[0].payload.header.channel_header.timestamp);
    console.log(transEvents[0].transactionData.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].collection_hashed_rwset)

    for (const transEvent of transEvents) {
        console.log(`*** transaction event: ${transEvent.transactionId}`);
        if (transEvent.privateData) {
            for (const namespace of transEvent.privateData.ns_pvt_rwset) {
                console.log(`    - private data: ${namespace.namespace}`);
                for (const collection of namespace.collection_pvt_rwset) {
                    console.log(`     - collection: ${collection.collection_name}`);

                    if (collection.rwset.reads) {
                        for (const read of collection.rwset.reads) {
                            console.log(read)
                           // console.log(`       - read set - ${BLUE}key:${RESET} ${read.key}  ${BLUE}value:${read.value.toString()}`);
                        }
                    }
                    if (collection.rwset.writes) {
                        for (const write of collection.rwset.writes) {
                            a = write.value.toString();  
                            b= JSON.parse(a);
                            console.log("##############",b)
                            let privateData = JSON.parse(b.MongoDBJson.replace(/\'/g, '"'));
                            console.log("*******************", privateData.orgName);
                            let validEvent = await checkBlockNumber(event.blockNumber.low)
                            if (validEvent) {
                                console.log("event is valid")
                                if (privateData.orgName != constants.orgName) {

                                    let returnResult=    request.post({
                                        url: 'http://13.235.241.13:3000/api/v1/blockchainEvents/blockchain_events',
                                        json: true,
                                        traditional: true,
                                        body: b
                                    },(err, result)=>{
                                        if(!err){
                                            console.log("Temp",temp)
                                            console.log("((((",result.body)
                                            temp++;
                                            return result.body
                                        }
                                        // console.log("Result Body========", result.body)
                                    })
                                } else {
                                    console.log("Skip the process of updating mongodb because organization is same",constants.orgName)
                                }
                            } else {
                                console.log("event is not valid")
                            }
                           
                            console.log("new",b);
                           
                            // console.log(`      - write set - ${BLUE}key:${RESET}${write.key} ${BLUE}is_delete:${RESET}${write.is_delete} ${BLUE}value:${RESET}${write.value.toString()}`);
                        }
                    }
                }
            }
         }
        }
    console.log("Block Event Received...")
    console.log("Block Number: ", event.blockNumber.low)
    console.log("Total Transactions: ", event.getTransactionEvents().length)

    console.log("Tx Type: ", event.blockData.data.data[0].payload.header.channel_header.type)
    console.log("Tx Timestamp: ", event.blockData.data.data[0].payload.header.channel_header.timestamp)

    console.log("Block Data[0].actions[0].ns_rwset[0]: ", event.blockData.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[0].namespace)
    console.log("Block Data[0].actions[0].ns_rwset[1]: ", event.blockData.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].namespace)
    console.log("Block Data[0].actions[0].ns_rwset[1].writeset[0].KEY: ", event.blockData.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0].key)
   // const event_data=event.blockData.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0].value.toString();
    console.log("Block Data[0].actions[0].ns_rwset[1].writeset[0].VALUE: ", event.blockData.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0].value.toString())
    
    // request.post({
    //     url: 'http://localhost:3000/api/v1/user/asset_event',
    //     json: true,
    //     traditional: true,
    //     body: transEvents,
    //     headers: {
           
    //     }
    // },(err, result)=>{
    //    // console.log("Result Body========", result.body)
    // })
}
var previousText = null;
function checkBlockNumber(number) {
  var currentText = number;
  if (previousText == currentText) {
      console.log("!!!!!!!!!!!",previousText,currentText)
      previousText = currentText;
    return false;
  } else {
    console.log("!!!!!!!!!!!",previousText,currentText)
    previousText = currentText;
     return true;
  }
}

const contractListener = async(event) =>{
    console.log("contract Event Received...")

    //console.log("Chaincode ID: ", chaincodeID)
    console.log("Event Name: ", event.eventName)
    console.log("Payload: ", event.payload.toString())
}

const commitListener = async(error, event) =>{
    console.log("Commit Listener Called...")
    if (error) {
        console.log("Error Occured in commit Listener: ", error)
    } else {
        console.log("Transaction Validity: ", event.peer.name)
        console.log("Contract Events Count: ", event.getContractEvents().length)
        const contractEvents = event.getContractEvents()
        console.log(contractEvents)
        console.log("Event Name: ", contractEvents[0].eventName)
        const txEvents = contractEvents.getTransactionEvents()
        console.log(txEvents)
        console.log(txEvents.isValid)
    }

}


module.exports = {
    contractListener,
    blockListener,
    commitListener
}