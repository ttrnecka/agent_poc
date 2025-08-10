db = db.getSiblingDB("poc");

db.policies.insertMany([
  {
    name: "brocade_cli",
    description: "Brocade CLI",
    file_name: "brocade_cli",
    versions: ["1.0.0", "1.0.1", "1.0.2", "1.0.3", "1.0.4"]
  },
  {
    name: "hpe3par_cli",
    description: "HPE 3PAR/Primera CLI",
    file_name: "hpe3par_cli",
    versions: ["1.0.0", "1.0.1", "2.0.0"]
  }
]);
db.collectors.insertMany([
  { name: "collector1", status: "OFFLINE" },
  { name: "collector2", status: "OFFLINE" },
]);

const collector1Id = db.collectors.findOne({ name: "collector1" })._id;

db.probes.insertMany([
    {
      collector_id: collector1Id,
      policy:   "brocade_cli",
      version:  "1.0.2",
      address:  "switch1",
      port:     2222,
      user:     "e",
      password: "secret"
    }
])
