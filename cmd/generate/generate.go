package main

import (
	"github.com/Fonzeca/Trackin/db"
	"gorm.io/gen"
)

func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./db/query",
		ModelPkgPath: "./db/model",
		Mode:         gen.WithoutContext,
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		/* FieldCoverable: true,*/
		// if you want generate field with unsigned integer type, set FieldSignable true
		/* FieldSignable: true,*/
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		WithUnitTest: false,
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	db.InitDB()
	defer db.CloseDB()
	g.UseDB(db.DB)

	logModel := g.GenerateModel("log")
	zonaModel := g.GenerateModel("zona")
	zonaVehiculoModel := g.GenerateModel("zona_vehiculos")

	g.ApplyBasic(logModel, zonaModel, zonaVehiculoModel)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	// g.ApplyBasic(model.User{}, g.GenerateModel("company"), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))

	// apply diy interfaces on structs or table models

	// execute the action of code generation
	g.Execute()
}
