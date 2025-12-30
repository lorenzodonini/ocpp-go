import {randomIntBetween, randomString} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import ws from 'k6/ws';
import {check, sleep} from 'k6';

// OCPP 2.1 Core profile message structure constants
const MESSAGE_TYPE = {
    CALL: 2,
    CALL_RESULT: 3,
    CALL_ERROR: 4

};

const ACTIONS = {
    BOOT_NOTIFICATION: 'BootNotification',
    HEARTBEAT: 'Heartbeat',
    STATUS_NOTIFICATION: 'StatusNotification',
    AUTHORIZE: 'Authorize',
    TRANSACTION_EVENT: 'TransactionEvent',
    METER_VALUES: 'MeterValues'
};

// OCPP 2.1 Core profile specific constants
const BOOT_REASONS = [
    'ApplicationReset',
    'FirmwareUpdate',
    'LocalReset',
    'PowerUp',
    'RemoteReset',
    'ScheduledReset',
    'Triggered',
    'Unknown',
    'Watchdog'
];

const CONNECTOR_STATUSES = [
    'Available',
    'Occupied',
    'Reserved',
    'Unavailable',
    'Faulted'
];

// Core profile incoming request actions from central system
const INCOMING_ACTIONS = {
    CHANGE_AVAILABILITY: 'ChangeAvailability',
    CHANGE_CONFIGURATION: 'ChangeConfiguration',
    CLEAR_CACHE: 'ClearCache',
    DATA_TRANSFER: 'DataTransfer',
    GET_CONFIGURATION: 'GetConfiguration',
    RESET: 'Reset',
    UNLOCK_CONNECTOR: 'UnlockConnector',
    REMOTE_START_TRANSACTION: 'RemoteStartTransaction',
    REMOTE_STOP_TRANSACTION: 'RemoteStopTransaction'
};

// Generate unique charging station ID
function generateChargingStationId() {
    return `CS_${__VU}_${randomString(8)}`;
}

// Generate OCPP 2.1 Core profile BootNotification message
function createBootNotification(chargingStationId: string) {
    const messageId = randomString(8);
    const payload = {
        reason: BOOT_REASONS[randomIntBetween(0, BOOT_REASONS.length - 1)],
        chargingStation: {
            serialNumber: chargingStationId,
            model: `Model_${randomString(5)}`,
            vendorName: `Vendor_${randomString(5)}`,
            firmwareVersion: `v${randomIntBetween(1, 9)}.${randomIntBetween(0, 9)}.${randomIntBetween(0, 9)}`,
            modem: {
                iccid: randomString(20),
                imsi: randomString(15)
            }
        }
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.BOOT_NOTIFICATION, payload]);
}

// Generate OCPP 2.1 Core profile Heartbeat message
function createHeartbeat() {
    const messageId = randomString(8);
    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.HEARTBEAT, {}]);
}

// Generate OCPP 2.1 Core profile StatusNotification message
function createStatusNotification(chargingStationId: string, evseId: number, connectorId: number, status: string) {
    const messageId = randomString(8);
    const payload = {
        evse: {
            id: evseId,
            connectorId: connectorId
        },
        timestamp: new Date().toISOString(),
        connectorStatus: status
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.STATUS_NOTIFICATION, payload]);
}

// Generate OCPP 2.1 Core profile Authorize message
function createAuthorize(idToken: string) {
    const messageId = randomString(8);
    const payload = {
        idToken: {
            idToken: idToken,
            type: 'ISO14443'
        }
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.AUTHORIZE, payload]);
}

// Generate OCPP 2.1 Core profile TransactionEvent message
function createTransactionEvent(chargingStationId: string, evseId: number, transactionId: string, eventType: string) {
    const messageId = randomString(8);
    const payload = {
        eventType: eventType,
        timestamp: new Date().toISOString(),
        triggerReason: 'Authorized',
        seqNo: randomIntBetween(1, 1000),
        transactionInfo: {
            transactionId: transactionId,
            chargingState: eventType === 'Started' ? 'Charging' : 'Idle'
        },
        evse: {
            id: evseId,
            connectorId: 1
        }
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.TRANSACTION_EVENT, payload]);
}

// Generate OCPP 2.1 Core profile MeterValues message
function createMeterValues(evseId: number, transactionId?: string) {
    const messageId = randomString(8);
    const payload = {
        evseId: evseId,
        transactionId: transactionId,
        meterValue: [{
            timestamp: new Date().toISOString(),
            sampledValue: [{
                value: randomIntBetween(0, 100000),
                context: 'Sample.Periodic',
                measurand: 'Energy.Active.Import.Register',
                unitOfMeasure: {
                    unit: 'Wh',
                    multiplier: 0
                }
            }]
        }]
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.METER_VALUES, payload]);
}

// Generate OCPP 2.1 Core profile CallResult response
function createCallResult(messageId: string, payload: any) {
    return JSON.stringify([MESSAGE_TYPE.CALL_RESULT, messageId, payload]);
}

// Generate OCPP 2.1 Core profile CallError response
function createCallError(messageId: string, errorCode: string, errorDescription: string) {
    return JSON.stringify([MESSAGE_TYPE.CALL_ERROR, messageId, errorCode, errorDescription, {}]);
}

// Handle incoming Core profile ChangeAvailability request
function handleChangeAvailability(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ChangeConfiguration request
function handleChangeConfiguration(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ClearCache request
function handleClearCache(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile DataTransfer request
function handleDataTransfer(request: any) {
    return {
        status: 'Accepted',
        data: `Response to ${request.vendorId}:${request.messageId || 'default'}`
    };
}

// Handle incoming Core profile GetConfiguration request
function handleGetConfiguration(request: any) {
    return {
        configurationKey: [
            {
                key: 'HeartbeatInterval',
                readonly: false,
                value: '60'
            },
            {
                key: 'ConnectionTimeOut',
                readonly: false,
                value: '60'
            }
        ],
        unknownKey: []
    };
}

// Handle incoming Core profile Reset request
function handleReset(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile UnlockConnector request
function handleUnlockConnector(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile RemoteStartTransaction request
function handleRemoteStartTransaction(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile RemoteStopTransaction request
function handleRemoteStopTransaction(request: any) {
    return {
        status: 'Accepted'
    };
}

// Route incoming Core profile requests to appropriate handlers
function handleIncomingRequest(action: string, payload: any, messageId: string) {
    let response;

    try {
        switch (action) {
            case INCOMING_ACTIONS.CHANGE_AVAILABILITY:
                response = handleChangeAvailability(payload);
                break;
            case INCOMING_ACTIONS.CHANGE_CONFIGURATION:
                response = handleChangeConfiguration(payload);
                break;
            case INCOMING_ACTIONS.CLEAR_CACHE:
                response = handleClearCache(payload);
                break;
            case INCOMING_ACTIONS.DATA_TRANSFER:
                response = handleDataTransfer(payload);
                break;
            case INCOMING_ACTIONS.GET_CONFIGURATION:
                response = handleGetConfiguration(payload);
                break;
            case INCOMING_ACTIONS.RESET:
                response = handleReset(payload);
                break;
            case INCOMING_ACTIONS.UNLOCK_CONNECTOR:
                response = handleUnlockConnector(payload);
                break;
            case INCOMING_ACTIONS.REMOTE_START_TRANSACTION:
                response = handleRemoteStartTransaction(payload);
                break;
            case INCOMING_ACTIONS.REMOTE_STOP_TRANSACTION:
                response = handleRemoteStopTransaction(payload);
                break;
            default:
                console.log(`VU ${__VU}: Unknown incoming action: ${action}`);
                return createCallError(messageId, 'NotImplemented', `Action ${action} not implemented`);
        }

        return createCallResult(messageId, response);
    } catch (error) {
        console.log(`VU ${__VU}: Error handling ${action}: ${error}`);
        return createCallError(messageId, 'InternalError', `Failed to process ${action}`);
    }
}

export default function () {
    const chargingStationId = generateChargingStationId();

    // Configurable WebSocket URL from environment variables
    const wsHost = __ENV.WS_HOST || 'central-system';
    const wsPort = __ENV.WS_PORT || '8887';
    const wsProtocol = __ENV.WS_PROTOCOL || 'ws';
    const wsPath = __ENV.WS_PATH || '';

    const url = `${wsProtocol}://${wsHost}:${wsPort}${wsPath}/${chargingStationId}`;

    console.log(`VU ${__VU}: Connecting to ${url}`);

    const params = {
        // OCPP 2.1 requires specific WebSocket subprotocol
        headers: {
            'Sec-WebSocket-Protocol': 'ocpp2.1'
        }
    };

    // Track connection metrics
    const startTime = Date.now();
    let messageCount = 0;
    let reconnectCount = 0;
    const maxReconnects = randomIntBetween(3, 8); // Random number of reconnections per VU

    // Main connection loop with rapid connect/disconnect
    for (let i = 0; i < maxReconnects; i++) {
        const res = ws.connect(url, params, function (socket) {
            const connectionStart = Date.now();

            socket.on('open', function open() {
                const connectTime = Date.now() - connectionStart;
                console.log(`VU ${__VU}: Connected in ${connectTime}ms (attempt ${i + 1}/${maxReconnects})`);

                // Send BootNotification immediately after connection
                const bootMsg = createBootNotification(chargingStationId);
                socket.send(bootMsg);
                messageCount++;

                // Disconnect after random time (100-200ms)
                const disconnectDelay = randomIntBetween(100, 200);
                setTimeout(() => {
                    console.log(`VU ${__VU}: Disconnecting after ${disconnectDelay}ms`);
                    socket.close();
                }, disconnectDelay);
            });

            socket.on('message', function (message) {
                try {
                    const data = JSON.parse(message);

                    if (data[0] === MESSAGE_TYPE.CALL_RESULT) {
                        console.log(`VU ${__VU}: Received response for message ${data[1]}`);
                    } else if (data[0] === MESSAGE_TYPE.CALL_ERROR) {
                        console.log(`VU ${__VU}: Received error for message ${data[1]}: ${data[2]}`);
                    } else if (data[0] === MESSAGE_TYPE.CALL) {
                        // Handle incoming request from central system
                        const messageId = data[1];
                        const action = data[2];
                        const payload = data[3] || {};

                        console.log(`VU ${__VU}: Received incoming request: ${action} (ID: ${messageId})`);

                        // Generate and send response
                        const response = handleIncomingRequest(action, payload, messageId);
                        socket.send(response);

                        // Track incoming request handling
                        messageCount++;
                    }
                } catch (e) {
                    console.log(`VU ${__VU}: Received non-JSON message: ${message}`);
                }
            });

            socket.on('close', function () {
                const totalTime = Date.now() - connectionStart;
                console.log(`VU ${__VU}: Disconnected after ${totalTime}ms`);
                reconnectCount++;
            });

            socket.on('error', function (error) {
                console.log(`VU ${__VU}: WebSocket error: ${error}`);
            });
        });

        // Verify connection was successful
        check(res, {
            'Connected successfully': (r) => r && r.status === 101,
            'Connection time < 200ms': (r) => r && (Date.now() - startTime) < 200
        });

        // Small delay between reconnections to avoid overwhelming the server
        sleep(randomIntBetween(50, 100) / 1000);
    }

    // Final metrics
    console.log(`VU ${__VU}: Completed ${reconnectCount} connections with ${messageCount} messages in ${Date.now() - startTime}ms`);
}
