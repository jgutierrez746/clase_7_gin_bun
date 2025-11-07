package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

// DB es la instancia global de Bun para MySQL
var DB *bun.DB

// InitDB inicializa la conexión a MySQL usando Bun.
// ds: Cadena de conexión, e: "user:pass@tcp(localhost:3306)/dbname?parseTime=true".
func InitDB(dsn string) error {
	// Abre la conexión con el driver estándar de MySQL.
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error abriendo conexión MySQL: %w", err)
	}

	// Crea la instancia de Bun con el dialecto MySQL.
	db := bun.NewDB(sqldb, mysqldialect.New())

	// Prueba de conexión.
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error conectando a MySQL: %w", err)
	}

	log.Println("Conexión a MySQL exitosa con Bun.")
	DB = db
	return nil
}

// SelectAll realiza un SELECT de todas las filas de una tabla y las escanea en un slice de structs.
// Ej: var users []User; err := SelectAll(ctx, "users", &users)
func SelectAll(ctx context.Context, tableName string, dest interface{}) error {
	if DB == nil {
		return fmt.Errorf("db no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}

	// Inferir nombre de talba si no se proporciona
	if tableName == "" {
		tableName = inferTableName(dest)
		if tableName == "" {
			return fmt.Errorf("nombre de tabla no proporcionado ni inferido del tag bun:table")
		}
	}

	return DB.NewSelect().Table(tableName).Scan(ctx, dest)
}

// SelectOne realiza un SELECT de una sola fila de una tabla con clausula WHERE.
// Ej: var users User; err := SelectOne(ctx, "users", &user, "id = ?", 1)
func SelectOne(ctx context.Context, tableName string, dest interface{}, where string, args ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}

	// Inferir nombre de talba si no se proporciona
	if tableName == "" {
		tableName = inferTableName(dest)
		if tableName == "" {
			return fmt.Errorf("nombre de tabla no proporcionado ni inferido del tag bun:table")
		}
	}

	q := DB.NewSelect().Table(tableName).Where(where, args...)
	return q.Scan(ctx, dest)
}

// Update actualiza filas en una tabla usando un modelo (struct) y cláusula WHERE.
// Retorna el número de filas afectadas (int64).
// Ej: user.Name = "Nuevo"; affected, err := Update(ctx, "users", &user, "id = ?", user.ID)
func Update(ctx context.Context, table string, model interface{}, where string, args ...interface{}) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewUpdate().Table(table).Model(model).Where(where, args...)
	res, err := q.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete borra filas de una tabla con cláusula WHERE.
// Retorna el número de filas afectadas (int64).
// Ej: affected, err := Delete(ctx, "users", "id = ?", 1)
func Delete(ctx context.Context, table string, where string, args ...interface{}) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewDelete().Table(table).Where(where, args...)
	res, err := q.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// SelectWithJoin realiza un SELECT con JOINs en una tabla principal, escaneando en un slice de structs.
// Ej:
//   - joins: []string{"LEFT JOIN orders ON orders.user_id = users.id"}
//   - where: "users.id > ?", 10
//
// Nota: El modelo (dest) debe mapear todos los campos de las tablas unidas.
func SelectConJoin(ctx context.Context, mainTable string, joins []string, dest interface{}, where string, args ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewSelect().Table(mainTable)

	// Agregar JOINs dinámicamente.
	for _, join := range joins {
		q = q.Join(join)
	}

	// WHERE opcional.
	if where != "" {
		q = q.Where(where, args...)
	}
	return q.Scan(ctx, dest)
}

// Insert inserta un modelo (struct) en la tabla (infiriendo del tag bun:table)
// Bun maneja autoincrement (ID)
func Insert(ctx context.Context, model interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}

	_, err := DB.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error insertando: %w", err)
	}
	log.Printf("Registro insertado exitosamente en tabla inferida del modelo.")
	return nil
}

// InsertBatch inserta múltiples modelos en batch (más eficiente para muchos registros)
func InsertBatch(ctx context.Context, models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}

	_, err := DB.NewInsert().Model(&models).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error insertando batch: %w", err)
	}
	log.Printf("Batch de %d registros insertado exitosamente.", len(models))
	return nil
}

func CreateTable(ctx context.Context, tableName string, model interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}

	// Inferir nombre de talba si no se proporciona
	if tableName == "" {
		tableName = inferTableName(model)
		if tableName == "" {
			return fmt.Errorf("nombre de tabla no proporcionado ni inferido del tag bun:table")
		}
	}

	_, err := DB.NewCreateTable().Table(tableName).Model(model).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creando tabla %s: %w", tableName, err)
	}
	log.Printf("Tabla %s creada exitosamente con esquema del modelo.", tableName)
	return nil
}

func inferTableName(model interface{}) string {
	rt := reflect.TypeOf(model)
	if rt == nil {
		return ""
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// Si es slice, ve al tipo del elemento
	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("bun")
		if tag != "" && strings.Contains(tag, "table:") {
			parts := strings.Split(tag, "table:")[1]
			if idx := strings.Index(parts, ","); idx > 0 {
				return strings.TrimSpace(parts[:idx])
			}
			return strings.TrimSpace(parts)
		}
	}
	return ""
}
