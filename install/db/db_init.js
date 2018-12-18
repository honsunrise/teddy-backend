db = db.getSiblingDB('admin');
db.createUser(
    {
        user: "zhsyourai",
        pwd: "zhscoderc..",
        roles: ["clusterAdmin", "readAnyDatabase", "dbAdmin", "userAdmin"]
    }
);

db = db.getSiblingDB('teddy');
db.createUser(
    {
        user: "teddy",
        pwd: "teddy",
        roles: [{role: "readWrite", db: "teddy"}]
    }
);