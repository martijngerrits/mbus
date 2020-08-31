package mbus

var products = map[string /* Manufacturer */]map[byte /* Version */]string /* Product name */ {
    "LAS": {
        0x01: "LAN-WMBUS-E-VOC",
        0x03: "LAN-WMBUS-E-CO2",
        0x07: "LAN-WMBUS-C-TH / LAN-WMBUS-G2-TH",
        0x0B: "LAN-WMBUS-G2-LDS",
        0x14: "LAN-WMBUS-G2-DC / LAN-WMBUS-G2-P",
        0x1E: "LAN-WMBUS-G2-EXT / LAN-WMBUS-G2-OOP",
    },
}
