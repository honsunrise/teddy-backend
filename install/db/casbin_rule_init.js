db = db.getSiblingDB('teddy');

// for anonuser
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/tags", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/tags/:tagID", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info/:id", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info/:id/segment", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info/:id/segment/:segID", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info/:id/segment/:segID/value", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/info/:id/segment/:segID/value/:valID", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/favorite/info/:id", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/thumbUp/info/:id", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/thumbDown/info/:id", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "", v1: "/v1/anon/content/search", v2: "GET"});

// fot user group
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/tags", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/tags/:tagID", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info/:id", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info/:id/segment", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info/:id/segment/:segID", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info/:id/segment/:segID/value", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/info/:id/segment/:segID/value/:valID", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/favorite/info/:id", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/thumbUp/info/:id", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/thumbDown/info/:id", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/anon/content/search", v2: "GET"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id", v2: "DELETE"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment/:segID", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment/:segID", v2: "DELETE"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment/:segID/value", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment/:segID/value/:valID", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/info/:id/segment/:segID/value/:valID", v2: "DELETE"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/favorite/user", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/favorite/info/:id", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/favorite/info/:id", v2: "DELETE"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbUp/user", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbUp/info/:id", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbUp/info/:id", v2: "DELETE"});

db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbDown/user", v2: "GET"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbDown/info/:id", v2: "POST"});
db.casbin_rule.insert({ptype: "p", v0: "user", v1: "/v1/auth/content/thumbDown/info/:id", v2: "DELETE"});

// user groups
db.casbin_rule.insert({ptype: "g", v0: "6887778051", v1: "user"});
