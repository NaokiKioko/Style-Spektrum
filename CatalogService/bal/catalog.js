const { report } = require('process');
const Databaseinterface = require('../dal/mongoDB.js'); // Import the class
const { ObjectId } = require('mongodb');

async function GetCatalogs(dal) {
    return await dal.interface("get", "Catalog", {});
}

async function GetCatalog(dal, id) {
    return await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(id) });
}

async function GetCatalogbyTags(dal, tags) {
    return await dal.interface("get", "Catalog", { "Tags": { $in: tags } });
}


async function PostCatalog(dal, catalog) {
    let code = await dal.interface("post", "Catalog", catalog);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(dal, catalog.Tags);
    }
    return code;
}

async function PatchCatalog(dal, id, catalog) {
    if (catalog._id) {
        delete catalog._id;
    }
    code = await dal.interface("patch", "Catalog", [{ _id: ObjectId.createFromHexString(id) }, catalog]);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function DeleteCatalog(dal, id) {
    return await dal.interface("delete", "Catalog", { _id: ObjectId.createFromHexString(id) });
}

async function GetTags(dal) {
    return await dal.interface("get", "Tags", {});
}

async function PostTags(dal, tagNames) {
    for (let i = 0; i < tagNames.length; i++) {
        let tag = await dal.interface("get", "Tags", { "Name": tagNames[i] });
        if (tag === 500) {
            return 500;
        }
        if (tag.length === 0) {
            await dal.interface("post", "Tags", { "Name": tagNames[i], "Favoritecount": 0 });
        }
    }
}

class Report {
    constructor(ReportedID, Field, NewContent, ReporterEmail) {
        this.ReportedID = ReportedID;
        this.Field = Field;
        this.NewContent = NewContent;
        this.ReporterEmail = [ReporterEmail];
        this.Popularity = 1;
    }
}
async function PostReport(dal, report) {
    let create = true;
    let code = 201;
    let newReport = new Report(report.id, report.Field, report.NewContent, report.Email);
    let product = await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(newReport.ReportedID) });
    if (product === 500) {
        return 500;
    }
    let reports = await dal.interface("get", "Reports", { "ReportedID": newReport.ReportedID, "Field": newReport.Field });
    for (let i = 0; i < reports.length; i++) {
        if (reports[i].NewContent === newReport.NewContent && reports[i].ReporterEmail.indexOf(newReport.ReporterEmail[0]) === -1) {
            reports[i].Popularity++;
            reports[i].ReporterEmail.push(newReport.ReporterEmail[0]);
            code = await dal.interface("patch", "Reports", [{ _id: reports[i]._id }, reports[i]]);
            create = false;
            break;
        } else if (reports[i].NewContent === newReport.NewContent && reports[i].ReporterEmail.indexOf(newReport.ReporterEmail[0]) != -1) {
            return 201;
        }
    }
    if (create) {
        code = await dal.interface("post", "Reports", newReport);
        if (code === 500) {
            return 500;
        }
    }
    await EnactReports(dal, newReport.ReportedID, newReport.Field, null, newReport.NewContent);
    return code;
}

class ReportItemTag {
    constructor(id, tagName, ReporterEmail) {
        this.ReportedID = id;
        this.Field = "Tags";
        this.TagName = tagName;
        this.Popularity = 1;
        this.ReporterEmail = [ReporterEmail];
    }
}

async function PostReportItemTag(dal, report) {
    let code;
    let create = true;
    let newReport = new ReportItemTag(report.id, report.tagName, report.Email);
    let product = await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(report.id) });
    if (product === 500) {
        return 500;
    }
    let reports = await dal.interface("get", "Reports", { "ReportedID": newReport.ReportedID, "Field": "Tags" });
    for (let i = 0; i < reports.length; i++) {
        if (reports[i].TagName === newReport.TagName && reports[i].ReporterEmail.indexOf(newReport.ReporterEmail[0]) === -1) {
            reports[i].Popularity++;
            reports[i].ReporterEmail.push(newReport.ReporterEmail[0]);
            code = await dal.interface("patch", "Reports", [{ _id: reports[i]._id }, reports[i]]);
            create = false;
            break;
        }
        else if (reports[i].TagName === newReport.TagName && reports[i].ReporterEmail.indexOf(newReport.ReporterEmail[0]) != -1) {
            return 201;
        }
    }
    if (create) {
        code = await dal.interface("post", "Reports", newReport);
        if (code === 500) {
            return 500;
        }
    }
    await EnactReports(dal, newReport.ReportedID, "Tags", newReport.TagName);
    return code;
}

async function EnactReports(dal, ReportedID, Field, TagName, NewContent) {
    if (Field === "Tags" && TagName != null) {
        let reports = await dal.interface("get", "Reports", { "ReportedID": ReportedID, "Field": Field, "TagName": TagName });
        if (reports === 500) {
            return 500;
        }
        let report = reports[0];
        if (report.Popularity >= 5) {
            let catalog = await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(ReportedID) });
            if (catalog === 500) {
                return 500;
            }
            let tags = catalog[0].Tags;
            let index = tags.indexOf(TagName);
            if (index === -1) {
                tags.push(TagName);
            } else {
                tags.splice(index, 1);
            }
            await dal.interface("patch", "Catalog", [{ _id: ObjectId.createFromHexString(ReportedID) }, { Tags: tags }]);
            await DeleteReports(dal, ReportedID, Field);
        }
    } else {
        let reports = await dal.interface("get", "Reports", { "ReportedID": ReportedID, "Field": Field, "NewContent": NewContent });
        if (reports === 500) {
            return 500;
        }
        let report = reports[0];
        if (report.Popularity >= 5) {
            await dal.interface("patch", "Catalog", [{ _id: ObjectId.createFromHexString(ReportedID) }, { [Field]: NewContent }]);
            await DeleteReports(dal, ReportedID, Field);
        }
    }
    return 201;
}

async function DeleteReports(dal, ReportedID, Field) {
    let reports = await dal.interface("get", "Reports", { "ReportedID": ReportedID, "Field": Field });
    if (reports === 500) {
        return 500;
    }
    for (let i = 0; i < reports.length; i++) {
        await dal.interface("delete", "Reports", { _id: reports[i]._id });
    }
}

async function GetReports(dal, id) {
    return await dal.interface("get", "Reports", {ReportedID: id});
}

module.exports = {
    GetCatalogs,
    GetCatalog,
    GetCatalogbyTags,
    PostCatalog,
    PatchCatalog,
    DeleteCatalog,
    GetTags,
    PostReport,
    PostReportItemTag,
    GetReports
}