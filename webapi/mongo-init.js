db = db.getSiblingDB("poc");

db.policies.insertMany([
  {
    key: "brocade_cli",
    name: "Brocade CLI",
    file_name: "brocade_cli",
    versions: ["1.0.0", "1.0.1", "1.0.2", "1.0.3", "1.0.4"]
  },
  {
    key: "hpe3par_cli",
    name: "HPE 3PAR/Primera CLI",
    file_name: "hpe3par_cli",
    versions: ["1.0.0", "1.0.1", "2.0.0"]
  }
]);

db.probes.insertMany([{"id":"c473a9d6-9399-40c9-b5e4-1338c16eaebc","collector":"collector1","policy":"brocade_cli","version":"1.0.2","address":"switch1","port":"2222","user":"e","password":"secret"}]
)

db.collectors.insertMany([
  { key: "collector1", data: {} },
  { key: "collector2", data: {} }
]);