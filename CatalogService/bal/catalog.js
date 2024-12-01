const { report } = require('process');
const Databaseinterface = require('../dal/mongoDB.js'); // Import the class
const { ObjectId } = require('mongodb');

async function GetCatalogs(dal) {
    return await dal.interface("get", "Catalog", {});
}

async function GetCatalog(dal, id) {
    return await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(id) });
}

async function GetCatalogsByFeild(dal, name, value) {
    return await dal.interface("get", "Catalog", { [name]: value });
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
    if (!IsField(report.field)) {
        return 400;
    }
    let create = true;
    let code = 201;
    let newReport = new Report(report.id, report.field, report.NewContent, report.Email);
    let reports = await dal.interface("get", "Reports", { "ReportedID": report.id, "Field": report.field, "NewContent": report.NewContent });
    if (reports.length == 1) {
        create = false;
        let report = reports[0];
        let index = report.ReporterEmail.indexOf(newReport.ReporterEmail[0]);
        if (index === -1) {
            report.ReporterEmail.push(newReport.ReporterEmail[0]);
            report.Popularity++;
            await dal.interface("patch", "Reports", [{ _id: report._id }, { ReporterEmail: report.ReporterEmail, Popularity: report.Popularity }]);
            if (report.Popularity >= 5) {
                await EnactReports(dal, report);
            }
        }
    } else if (reports.length == 0) {
        if (create) {
            code = await dal.interface("post", "Reports", newReport);
            if (code === 500) {
                return 500;
            }
        }
    } else {
        for (let i = 0; i < reports.length; i++) {
            await DeleteReport(dal, reports[i]._id);
        }
        return 500;
    }
    return code;
}

async function EnactReports(dal, report) {
    let catalog = await dal.interface("get", "Catalog", { _id: ObjectId.createFromHexString(report.ReportedID) });
    if (catalog === 500) {
        return 500;
    }
    product = catalog[0];
    if (report.Field === "Tags" || report.Field === "Images") {
        let index = product[report.Field].indexOf(report.NewContent);
        if (index === -1) {
            product.Tags.push(report.NewContent);
        } else {
            product.Tags.splice(index, 1);
        }
        let code = await dal.interface("patch", "Catalog", [{ _id: ObjectId.createFromHexString(report.ReportedID) }, product]);
        if (code === 500) {
            return 500;
        }
    } else {
        product[report.Field] = report.NewContent;
        let code = await dal.interface("patch", "Catalog", [{ _id: ObjectId.createFromHexString(report.ReportedID) }, product]);
        if (code === 500) {
            return 500;
        }
    }
    let code = await DeleteReport(dal, report._id);
    return code;
}

async function DeleteReport(dal, id) {
    return await dal.interface("delete", "Reports", { _id: id });
}

async function GetReports(dal, id) {
    return await dal.interface("get", "Reports", { "ReportedID": id });
}

async function GetReportsByField(dal, id, field) {
    return await dal.interface("get", "Reports", { "ReportedID": id, "Field": field, });
}

async function DeleteReport(dal, id) {
    return await dal.interface("delete", "Reports", { _id: id });
}

function IsField(field) {
    let fields = ["Name", "Description", "URL", "Price", "Rating", "Tags", "Images"];
    for (let i = 0; i < fields.length; i++) {
        if (field === fields[i]) {
            return true;
        }
    }
    return false;
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
    GetReports,
    GetReportsByField,
    GetCatalogsByFeild
}