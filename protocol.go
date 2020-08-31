package mbus

const (
    // Frame types
    FRAME_TYPE_ANY byte = 0x00
    FRAME_TYPE_ACK = 0x01
    FRAME_TYPE_SHORT = 0x02
    FRAME_TYPE_CONTROL = 0x03
    FRAME_TYPE_LONG = 0x04

    FRAME_ACK_BASE_SIZE = 1
    FRAME_SHORT_BASE_SIZE = 5
    FRAME_CONTROL_BASE_SIZE = 9
    FRAME_LONG_BASE_SIZE = 9

    FRAME_BASE_SIZE_ACK = 1
    FRAME_BASE_SIZE_SHORT = 5
    FRAME_BASE_SIZE_CONTROL = 9
    FRAME_BASE_SIZE_LONG = 9

    FRAME_FIXED_SIZE_ACK = 1
    FRAME_FIXED_SIZE_SHORT = 5
    FRAME_FIXED_SIZE_CONTROL = 6
    FRAME_FIXED_SIZE_LONG = 6

    // Frame start/stop bits
    FRAME_ACK_START = 0xE5
    FRAME_SHORT_START = 0x10
    FRAME_CONTROL_START = 0x68
    FRAME_LONG_START = 0x68
    FRAME_STOP = 0x16

    //
    // Control field
    //
    CONTROL_FIELD_DIRECTION = 0x07
    CONTROL_FIELD_FCB = 0x06
    CONTROL_FIELD_ACD = 0x06
    CONTROL_FIELD_FCV = 0x05
    CONTROL_FIELD_DFC = 0x05
    CONTROL_FIELD_F3 = 0x04
    CONTROL_FIELD_F2 = 0x03
    CONTROL_FIELD_F1 = 0x02
    CONTROL_FIELD_F0 = 0x01

    CONTROL_MASK_SND_NKE = 0x40
    CONTROL_MASK_SND_NR = 0x44
    CONTROL_MASK_SND_UD = 0x53
    CONTROL_MASK_REQ_UD2 = 0x5B
    CONTROL_MASK_REQ_UD1 = 0x5A
    CONTROL_MASK_RSP_UD = 0x08

    CONTROL_MASK_FCB = 0x20
    CONTROL_MASK_FCV = 0x10

    CONTROL_MASK_ACD = 0x20
    CONTROL_MASK_DFC = 0x10

    CONTROL_MASK_DIR = 0x40
    CONTROL_MASK_DIR_M2S = 0x40
    CONTROL_MASK_DIR_S2M = 0x00

    //
    // Address field
    //
    ADDRESS_BROADCAST_REPLY = 0xFE
    ADDRESS_BROADCAST_NOREPLY = 0xFF
    ADDRESS_NETWORK_LAYER = 0xFD

    //
    // Control Information field
    //
    //Mode 1 Mode 2                   Application                   Definition in
    // 51h    55h                       data send                    EN1434-3
    // 52h    56h                  selection of slaves           Usergroup July  ́93
    // 50h                          application reset           Usergroup March  ́94
    // 54h                          synronize action                 suggestion
    // B8h                     set baudrate to 300 baud          Usergroup July  ́93
    // B9h                     set baudrate to 600 baud          Usergroup July  ́93
    // BAh                    set baudrate to 1200 baud          Usergroup July  ́93
    // BBh                    set baudrate to 2400 baud          Usergroup July  ́93
    // BCh                    set baudrate to 4800 baud          Usergroup July  ́93
    // BDh                    set baudrate to 9600 baud          Usergroup July  ́93
    // BEh                   set baudrate to 19200 baud              suggestion
    // BFh                   set baudrate to 38400 baud              suggestion
    // B1h           request readout of complete RAM content     Techem suggestion
    // B2h          send user data (not standardized RAM write) Techem suggestion
    // B3h                 initialize test calibration mode      Usergroup July  ́93
    // B4h                           EEPROM read                 Techem suggestion
    // B6h                         start software test           Techem suggestion
    // 90h to 97h              codes used for hashing           longer recommended
    CONTROL_INFO_DATA_SEND = 0x51
    CONTROL_INFO_DATA_SEND_MSB = 0x55
    CONTROL_INFO_SELECT_SLAVE = 0x52
    CONTROL_INFO_SELECT_SLAVE_MSB = 0x56
    CONTROL_INFO_APPLICATION_RESET = 0x50
    CONTROL_INFO_SYNC_ACTION = 0x54
    CONTROL_INFO_SET_BAUDRATE_300 = 0xB8
    CONTROL_INFO_SET_BAUDRATE_600 = 0xB9
    CONTROL_INFO_SET_BAUDRATE_1200 = 0xBA
    CONTROL_INFO_SET_BAUDRATE_2400 = 0xBB
    CONTROL_INFO_SET_BAUDRATE_4800 = 0xBC
    CONTROL_INFO_SET_BAUDRATE_9600 = 0xBD
    CONTROL_INFO_SET_BAUDRATE_19200 = 0xBE
    CONTROL_INFO_SET_BAUDRATE_38400 = 0xBF
    CONTROL_INFO_REQUEST_RAM_READ = 0xB1
    CONTROL_INFO_SEND_USER_DATA = 0xB2
    CONTROL_INFO_INIT_TEST_CALIB = 0xB3
    CONTROL_INFO_EEPROM_READ = 0xB4
    CONTROL_INFO_SW_TEST_START = 0xB6

    //Mode 1 Mode 2                   Application                   Definition in
    // 70h             report of general application errors     Usergroup March 94
    // 71h                      report of alarm status          Usergroup March 94
    // 72h   76h                variable data respond                EN1434-3
    // 73h   77h                 fixed data respond                  EN1434-3
    CONTROL_INFO_ERROR_GENERAL = 0x70
    CONTROL_INFO_STATUS_ALARM = 0x71

    CONTROL_INFO_RESP_FIXED = 0x73
    CONTROL_INFO_RESP_FIXED_MSB = 0x77

    CONTROL_INFO_RESP_VARIABLE = 0x72
    CONTROL_INFO_RESP_VARIABLE_MSB = 0x76

    //
    // data record fields
    //
    DATA_RECORD_DIF_MASK_INST = 0x00
    DATA_RECORD_DIF_MASK_MIN = 0x10
    DATA_RECORD_DIF_MASK_TYPE_INT32 = 0x04
    DATA_RECORD_DIF_MASK_DATA = 0x0F
    DATA_RECORD_DIF_MASK_FUNCTION = 0x30
    DATA_RECORD_DIF_MASK_STORAGE_NO = 0x40
    DATA_RECORD_DIF_MASK_EXTENTION = 0x80
    DATA_RECORD_DIF_MASK_NON_DATA = 0xF0
    DATA_RECORD_DIFE_MASK_STORAGE_NO = 0x0F
    DATA_RECORD_DIFE_MASK_TARIFF = 0x30
    DATA_RECORD_DIFE_MASK_DEVICE = 0x40
    DATA_RECORD_DIFE_MASK_EXTENSION = 0x80

    //
    // GENERAL APPLICATION ERRORS
    //
    ERROR_DATA_UNSPECIFIED = 0x00
    ERROR_DATA_UNIMPLEMENTED_CI = 0x01
    ERROR_DATA_BUFFER_TOO_LONG = 0x02
    ERROR_DATA_TOO_MANY_RECORDS = 0x03
    ERROR_DATA_PREMATURE_END = 0x04
    ERROR_DATA_TOO_MANY_DIFES = 0x05
    ERROR_DATA_TOO_MANY_VIFES = 0x06
    ERROR_DATA_RESERVED = 0x07
    ERROR_DATA_APPLICATION_BUSY = 0x08
    ERROR_DATA_TOO_MANY_READOUTS = 0x09

    //
    // FIXED DATA FLAGS
    //

    //
    // VARIABLE DATA FLAGS
    //
    VARIABLE_DATA_MEDIUM_OTHER = 0x00
    VARIABLE_DATA_MEDIUM_OIL = 0x01
    VARIABLE_DATA_MEDIUM_ELECTRICITY = 0x02
    VARIABLE_DATA_MEDIUM_GAS = 0x03
    VARIABLE_DATA_MEDIUM_HEAT_OUT = 0x04
    VARIABLE_DATA_MEDIUM_STEAM = 0x05
    VARIABLE_DATA_MEDIUM_HOT_WATER = 0x06
    VARIABLE_DATA_MEDIUM_WATER = 0x07
    VARIABLE_DATA_MEDIUM_HEAT_COST = 0x08
    VARIABLE_DATA_MEDIUM_COMPR_AIR = 0x09
    VARIABLE_DATA_MEDIUM_COOL_OUT = 0x0A
    VARIABLE_DATA_MEDIUM_COOL_IN = 0x0B
    VARIABLE_DATA_MEDIUM_HEAT_IN = 0x0C
    VARIABLE_DATA_MEDIUM_HEAT_COOL = 0x0D
    VARIABLE_DATA_MEDIUM_BUS = 0x0E
    VARIABLE_DATA_MEDIUM_UNKNOWN = 0x0F
    VARIABLE_DATA_MEDIUM_IRRIGATION = 0x10
    VARIABLE_DATA_MEDIUM_WATER_LOGGER = 0x11
    VARIABLE_DATA_MEDIUM_GAS_LOGGER = 0x12
    VARIABLE_DATA_MEDIUM_GAS_CONV = 0x13
    VARIABLE_DATA_MEDIUM_COLORIFIC = 0x14
    VARIABLE_DATA_MEDIUM_BOIL_WATER = 0x15
    VARIABLE_DATA_MEDIUM_COLD_WATER = 0x16
    VARIABLE_DATA_MEDIUM_DUAL_WATER = 0x17
    VARIABLE_DATA_MEDIUM_PRESSURE = 0x18
    VARIABLE_DATA_MEDIUM_ADC = 0x19
    VARIABLE_DATA_MEDIUM_SMOKE = 0x1A
    VARIABLE_DATA_MEDIUM_ROOM_SENSOR = 0x1B
    VARIABLE_DATA_MEDIUM_GAS_DETECTOR = 0x1C
    VARIABLE_DATA_MEDIUM_BREAKER_E = 0x20
    VARIABLE_DATA_MEDIUM_VALVE = 0x21
    VARIABLE_DATA_MEDIUM_CUSTOMER_UNIT = 0x25
    VARIABLE_DATA_MEDIUM_WASTE_WATER = 0x28
    VARIABLE_DATA_MEDIUM_GARBAGE = 0x29
    VARIABLE_DATA_MEDIUM_VOC = 0x2B
    VARIABLE_DATA_MEDIUM_SERVICE_UNIT = 0x30
    VARIABLE_DATA_MEDIUM_RC_SYSTEM = 0x36
    VARIABLE_DATA_MEDIUM_RC_METER = 0x37

    //
    // ABSTRACT DATA FORMAT (error, fixed or variable length)
    //
    DATA_TYPE_FIXED = 1
    DATA_TYPE_VARIABLE = 2
    DATA_TYPE_ERROR = 3

    DIB_DIF_WITHOUT_EXTENSION = 0x7F
    DIB_DIF_EXTENSION_BIT = 0x80
    DIB_VIF_WITHOUT_EXTENSION = 0x7F
    DIB_VIF_EXTENSION_BIT = 0x80
    DIB_DIF_MANUFACTURER_SPECIFIC = 0x0F
    DIB_DIF_MORE_RECORDS_FOLLOW = 0x1F
    DIB_DIF_IDLE_FILLER = 0x2F

    // Control Mask
    FRAME_DATA_LENGTH = 252

    MAX_PRIMARY_SLAVES = 250

    DATA_VARIABLE_HEADER_LENGTH = 12
)
