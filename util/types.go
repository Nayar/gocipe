package util

const (
	// WidgetTypeCheckbox represents widget of type checkbox
	WidgetTypeCheckbox = "checkbox"

	// WidgetTypeDate represents widget of type date
	WidgetTypeDate = "date"

	// WidgetTypeNumber represents widget of type number
	WidgetTypeNumber = "number"

	// WidgetTypePassword represents widget of type password
	WidgetTypePassword = "password"

	// WidgetTypeSelect represents widget of type select
	WidgetTypeSelect = "select"

	// WidgetTypeSelectRel represents widget of type select-rel
	WidgetTypeSelectRel = "select-rel"

	// WidgetTypeTextArea represents widget of type textarea
	WidgetTypeTextArea = "textarea"

	// WidgetTypeTextField represents widget of type textfield
	WidgetTypeTextField = "textfield"

	// WidgetTypeTime represents widget of type time
	WidgetTypeTime = "time"

	// WidgetTypeToggle represents widget of type toggle
	WidgetTypeToggle = "toggle"

	// RelationshipTypeManyMany represents a relationship of type Many-Many
	RelationshipTypeManyMany = "many-many"

	// RelationshipTypeOneOne represents a relationship of type One-One
	RelationshipTypeOneOne = "one-one"

	// RelationshipTypeOneMany represents a relationship of type One-Many
	RelationshipTypeOneMany = "one-many"

	// RelationshipTypeManyOne represents a relationship of type Many-One
	RelationshipTypeManyOne = "many-one"

	// PrimaryKeySerial indicates primary key - autogenerated number
	PrimaryKeySerial = "serial"

	// PrimaryKeyUUID indicates primary key - autogenerated string
	PrimaryKeyUUID = "uuid"

	// PrimaryKeyInt indicates primary key - manually assigned number
	PrimaryKeyInt = "int"

	// PrimaryKeyString indicates primary key - manually assigned string
	PrimaryKeyString = "string"
)

// Recipe represents a recipe to generate a project
type Recipe struct {
	// Container indicates whether or not container should be generated
	Bootstrap BootstrapOpts `json:"bootstrap"`

	// HTTP indicates whether http server code should be generated
	HTTP HTTPOpts `json:"http"`

	// Schema describes options for Schema generation
	Schema SchemaOpts `json:"schema"`

	// Crud describes options for Crud generation
	Crud CrudOpts `json:"crud"`

	// Rest describes options for Rest generation
	Rest RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation
	Vuetify VuetifyOpts `json:"vuetify"`

	// Entities lists entities to be generated
	Entities []Entity `json:"entities"`
}

// HTTPOpts represents options for http function generation
type HTTPOpts struct {
	// Generate indicates whether or not to generate http serve function
	Generate bool `json:"generate"`

	// Port represents default port to run application
	Port string `json:"port"`
}

// BootstrapOpts represents options for bootstrap function generation
type BootstrapOpts struct {
	// Generate indicates whether or not to generate bootstrap
	Generate bool `json:"generate"`

	// NoDB indicates that database connection should not be generated (default false)
	NoDB bool `json:"no_db"`

	// Settings represent list of settings to load during bootstrap into main package
	Settings []BootstrapSetting `json:"settings"`
}

// BootstrapSetting represents a setting required by the application and loaded during bootstrap
type BootstrapSetting struct {
	// Name represents name of setting
	Name string `json:"name"`

	// Description gives information on the setting (useful to display errors if not found)
	Description string `json:"description"`

	// Public indicates if setting should be accessible from all packages
	Public bool `json:"public"`
}

// SchemaOpts represents options for schema generation
type SchemaOpts struct {
	// Create whether or not to generate CREATE TABLE
	Create bool `json:"create"`

	// Drop whether or not to generate DROP IF EXISTS before CREATE
	Drop bool `json:"drop"`

	// Aggregate whether or not to generate schema into single file instead of separate files
	Aggregate bool `json:"aggregate"`

	// Path indicates in which path to generate the schema sql file
	Path string `json:"path"`
}

// CrudOpts represents which crud functions should be generated
type CrudOpts struct {
	// Create indicates whether or not function for INSERT should be generated
	Create bool `json:"create"`

	// Read indicates whether or not function for single select by id - SELECT WHERE id = id should be generated
	Read bool `json:"read"`

	// ReadList indicates whether or not function for list select should be generated
	ReadList bool `json:"read_list"`

	// Update indicates whether or not function for UPDATE should be generated
	Update bool `json:"update"`

	// Delete indicates whether or not function for DELETE should be generated
	Delete bool `json:"delete"`

	// Merge indicates whether or not function for SQL Merge should be generated
	Merge bool `json:"merge"`

	// Hooks describes hooks options for CRUD generation
	Hooks CrudHooks `json:"hooks"`
}

// CrudHooks represents which crud hooks should be generated
type CrudHooks struct {

	// PreSave allows hook function to be executed before Save operation is performed
	PreSave bool `json:"pre_save"`

	// PostSave allows hook function to be executed after Save operation is performed
	PostSave bool `json:"post_save"`

	// PreRead allows hook function to be executed before Read operation is performed
	PreRead bool `json:"pre_read"`

	// PostRead allows hook function to be executed after Read operation is performed
	PostRead bool `json:"post_read"`

	// PreList allows hook function to be executed before List operation is performed
	PreList bool `json:"pre_list"`

	// PostList allows hook function to be executed after List operation is performed
	PostList bool `json:"post_list"`

	// PreDeleteSingle allows hook function to be executed before DeleteSingle operation is performed
	PreDeleteSingle bool `json:"pre_delete_single"`

	// PostDeleteSingle allows hook function to be executed after DeleteSingle operation is performed
	PostDeleteSingle bool `json:"post_delete_single"`

	// PreDeleteMany allows hook function to be executed before DeleteMany operation is performed
	PreDeleteMany bool `json:"pre_delete_many"`

	// PostDeleteMany allows hook function to be executed after DeleteMany operation is performed
	PostDeleteMany bool `json:"post_delete_many"`
}

// RestOpts represents which rest functions should be generated
type RestOpts struct {

	// Create indicates if http endpoint for POST method should be generated
	Create bool `json:"create"`

	// Read indicates if http endpoint for GET method (by id for single entity) should be generated
	Read bool `json:"read"`

	// ReadList indicates if http endpoint for GET method (by filters for many entities) should be generated
	ReadList bool `json:"read_list"`

	// Update indicates if http endpoint for PUT method should be generated
	Update bool `json:"update"`

	// Delete indicates if http endpoint for DELETE method should be generated
	Delete bool `json:"delete"`

	// Prefix indicates which prefix to use for routes
	Prefix string `json:"prefix"`

	// Hooks describes hooks options for REST generation
	Hooks RestHooks `json:"hooks"`
}

// RestHooks represents which rest hooks should be generated
type RestHooks struct {

	// PreCreate allows hook function to be executed before POST operations are done
	PreCreate bool `json:"pre_create"`

	// PostCreate allows hook function to be executed after POST operations are done
	PostCreate bool `json:"post_create"`

	// PreRead allows hook function to be executed before GET (single by id) operations are done
	PreRead bool `json:"pre_read"`

	// PostRead allows hook function to be executed after GET (single by id) operations are done
	PostRead bool `json:"post_read"`

	// PreList allows hook function to be executed before GET (many by filters) operations are done
	PreList bool `json:"pre_list"`

	// PostList allows hook function to be executed after GET (many by filters) operations are done
	PostList bool `json:"post_list"`

	// PreUpdate allows hook function to be executed before PUT operations are done
	PreUpdate bool `json:"pre_update"`

	// PostUpdate allows hook function to be executed after PUT operations are done
	PostUpdate bool `json:"post_update"`

	// PreDelete allows hook function to be executed before DELETE operations are done
	PreDelete bool `json:"pre_delete"`

	// PostDelete allows hook function to be executed after DELETE operations are done
	PostDelete bool `json:"post_delete"`
}

// VuetifyOpts represents options for vuetify generator
type VuetifyOpts struct {
	// Generate represents whether or not to generate vuetify assets
	Generate bool `json:"generate"`

	//Module represents the location where the gocipe module will be generated
	Module string `json:"module"`
}

// Entity represents a single entity to be generated
type Entity struct {
	// Name is the name of the entity
	Name string `json:"name"`

	// PrimaryKey indicates the nature of the primary key: serial (auto incremented number), uuid (auto generated string), int or string
	PrimaryKey string `json:"primary_key"`

	// Table is the name of the database table for the entity
	Table string `json:"table"`

	// TableConstraints represents an array of table constraints for the table definition
	TableConstraints []string `json:"table_constraints"`

	// Description is a description of the entity
	Description string `json:"description"`

	// Fields is a list of fields for the entity
	Fields []Field `json:"fields"`

	// Relationships represent relationship information between this entity and others
	Relationships []Relationship `json:"relationships"`

	// Schema describes options for Schema generation - overrides recipe level Schema config
	Schema *SchemaOpts `json:"schema"`

	// Crud describes options for Crud generation - overrides recipe level Crud config
	Crud *CrudOpts `json:"crud"`

	// Rest describes options for Rest generation - overrides recipe level Rest config
	Rest *RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation - overrides recipe level Vuetify config
	Vuetify *VuetifyOpts `json:"vuetify"`
}

// Field describes a field contained in an entity
type Field struct {
	// Label is the label for the field
	Label string `json:"label"`

	// Property represents code property information for the field
	Property FieldProperty `json:"property"`

	// Schema represents schema information for the field
	Schema FieldSchema `json:"schema"`

	// Widget represents widget information for the field
	Widget WidgetOpts `json:"widget"`

	// Filterable indicates if queries can be made using this field
	Filterable bool `json:"filterable"`
}

// FieldProperty represents code information for the field
type FieldProperty struct {
	// Name is the name of the property
	Name string `json:"name"`

	// Type is the data type of the property
	Type string `json:"type"`
}

// FieldSchema represents schema generation information for the field
type FieldSchema struct {
	// Field is the name of the field in database
	Field string `json:"field"`

	// Type is the data type for the field in database
	Type string `json:"type"`

	// Nullable indicates if null values are allowed in database for this field
	Nullable bool `json:"nullable"`

	// Default provides the default value for this field in database
	Default string `json:"default"`
}

// Relationship represents a relationship between this entity and another
type Relationship struct {
	// Entity is the name of the related entity
	Entity string `json:"entity"`

	// Type represents the type of relationship
	Type string `json:"type"`

	// Name represents the property name to be used for this relationship
	Name string `json:"name"`

	// JoinTable represents the other table in a many-many relationship
	JoinTable string `json:"join_table"`

	// ThisID represents the field in this entity used for the relationship
	ThisID string `json:"thisid"`

	// ThatID represents the field in the other entity used for the relationship
	ThatID string `json:"thatid"`
}

// WidgetOpts represents a UI widget
type WidgetOpts struct {
	// Type indicates which widget type is represented
	Type string `json:"type"`

	// Options represents options listed by this widget
	Options []WidgetOption `json:"options"`

	// Target represents a target endpoint to pull data for this widget
	Target WidgetTarget `json:"target"`

	// Multiple indicates that the field accepts multiple values
	Multiple bool
}

// WidgetOption represents an option for SelectRel widget type
type WidgetOption struct {
	// Value represents the stored value of the option
	Value string `json:"value"`
	// Label represents the displayed of the option
	Label string `json:"label"`
}

// WidgetTarget represents a target endpoint to pull data for this widget
type WidgetTarget struct {
	// Endpoint represents an endpoint to pull data from
	Endpoint string

	// Label which field to use for label on data endpoint
	Label string
}
