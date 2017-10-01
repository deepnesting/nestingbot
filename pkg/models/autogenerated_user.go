package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// ===== BEGIN of all query sets

// ===== BEGIN of query set UserQuerySet

// UserQuerySet is an queryset type for User
type UserQuerySet struct {
	db *gorm.DB
}

// NewUserQuerySet constructs new UserQuerySet
func NewUserQuerySet(db *gorm.DB) UserQuerySet {
	return UserQuerySet{
		db: db.Model(&User{}),
	}
}

func (qs UserQuerySet) w(db *gorm.DB) UserQuerySet {
	return NewUserQuerySet(db)
}

// All is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) All(ret *[]User) error {
	return qs.db.Find(ret).Error
}

// Count is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) Count() (int, error) {
	var count int
	err := qs.db.Count(&count).Error
	return count, err
}

// Create is an autogenerated method
// nolint: dupl
func (o *User) Create(db *gorm.DB) error {
	return db.Create(o).Error
}

// CreatedAtEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtEq(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at = ?", createdAt))
}

// CreatedAtGt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtGt(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at > ?", createdAt))
}

// CreatedAtGte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtGte(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at >= ?", createdAt))
}

// CreatedAtLt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtLt(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at < ?", createdAt))
}

// CreatedAtLte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtLte(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at <= ?", createdAt))
}

// CreatedAtNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) CreatedAtNe(createdAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("created_at != ?", createdAt))
}

// Delete is an autogenerated method
// nolint: dupl
func (o *User) Delete(db *gorm.DB) error {
	return db.Delete(o).Error
}

// Delete is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) Delete() error {
	return qs.db.Delete(User{}).Error
}

// DeletedAtEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtEq(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at = ?", deletedAt))
}

// DeletedAtGt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtGt(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at > ?", deletedAt))
}

// DeletedAtGte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtGte(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at >= ?", deletedAt))
}

// DeletedAtIsNotNull is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtIsNotNull() UserQuerySet {
	return qs.w(qs.db.Where("deleted_at IS NOT NULL"))
}

// DeletedAtIsNull is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtIsNull() UserQuerySet {
	return qs.w(qs.db.Where("deleted_at IS NULL"))
}

// DeletedAtLt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtLt(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at < ?", deletedAt))
}

// DeletedAtLte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtLte(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at <= ?", deletedAt))
}

// DeletedAtNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) DeletedAtNe(deletedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("deleted_at != ?", deletedAt))
}

// FirstNameEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) FirstNameEq(firstName string) UserQuerySet {
	return qs.w(qs.db.Where("first_name = ?", firstName))
}

// FirstNameIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) FirstNameIn(firstName string, firstNameRest ...string) UserQuerySet {
	iArgs := []interface{}{firstName}
	for _, arg := range firstNameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("first_name IN (?)", iArgs))
}

// FirstNameNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) FirstNameNe(firstName string) UserQuerySet {
	return qs.w(qs.db.Where("first_name != ?", firstName))
}

// FirstNameNotIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) FirstNameNotIn(firstName string, firstNameRest ...string) UserQuerySet {
	iArgs := []interface{}{firstName}
	for _, arg := range firstNameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("first_name NOT IN (?)", iArgs))
}

// GetUpdater is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) GetUpdater() UserUpdater {
	return NewUserUpdater(qs.db)
}

// IDEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDEq(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id = ?", ID))
}

// IDGt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDGt(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id > ?", ID))
}

// IDGte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDGte(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id >= ?", ID))
}

// IDIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDIn(ID uint, IDRest ...uint) UserQuerySet {
	iArgs := []interface{}{ID}
	for _, arg := range IDRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("id IN (?)", iArgs))
}

// IDLt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDLt(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id < ?", ID))
}

// IDLte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDLte(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id <= ?", ID))
}

// IDNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDNe(ID uint) UserQuerySet {
	return qs.w(qs.db.Where("id != ?", ID))
}

// IDNotIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) IDNotIn(ID uint, IDRest ...uint) UserQuerySet {
	iArgs := []interface{}{ID}
	for _, arg := range IDRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("id NOT IN (?)", iArgs))
}

// LastNameEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) LastNameEq(lastName string) UserQuerySet {
	return qs.w(qs.db.Where("last_name = ?", lastName))
}

// LastNameIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) LastNameIn(lastName string, lastNameRest ...string) UserQuerySet {
	iArgs := []interface{}{lastName}
	for _, arg := range lastNameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("last_name IN (?)", iArgs))
}

// LastNameNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) LastNameNe(lastName string) UserQuerySet {
	return qs.w(qs.db.Where("last_name != ?", lastName))
}

// LastNameNotIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) LastNameNotIn(lastName string, lastNameRest ...string) UserQuerySet {
	iArgs := []interface{}{lastName}
	for _, arg := range lastNameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("last_name NOT IN (?)", iArgs))
}

// Limit is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) Limit(limit int) UserQuerySet {
	return qs.w(qs.db.Limit(limit))
}

// One is used to retrieve one result. It returns gorm.ErrRecordNotFound
// if nothing was fetched
func (qs UserQuerySet) One(ret *User) error {
	return qs.db.First(ret).Error
}

// OrderAscByCreatedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderAscByCreatedAt() UserQuerySet {
	return qs.w(qs.db.Order("created_at ASC"))
}

// OrderAscByDeletedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderAscByDeletedAt() UserQuerySet {
	return qs.w(qs.db.Order("deleted_at ASC"))
}

// OrderAscByID is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderAscByID() UserQuerySet {
	return qs.w(qs.db.Order("id ASC"))
}

// OrderAscByTelegramID is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderAscByTelegramID() UserQuerySet {
	return qs.w(qs.db.Order("telegram_id ASC"))
}

// OrderAscByUpdatedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderAscByUpdatedAt() UserQuerySet {
	return qs.w(qs.db.Order("updated_at ASC"))
}

// OrderDescByCreatedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderDescByCreatedAt() UserQuerySet {
	return qs.w(qs.db.Order("created_at DESC"))
}

// OrderDescByDeletedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderDescByDeletedAt() UserQuerySet {
	return qs.w(qs.db.Order("deleted_at DESC"))
}

// OrderDescByID is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderDescByID() UserQuerySet {
	return qs.w(qs.db.Order("id DESC"))
}

// OrderDescByTelegramID is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderDescByTelegramID() UserQuerySet {
	return qs.w(qs.db.Order("telegram_id DESC"))
}

// OrderDescByUpdatedAt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) OrderDescByUpdatedAt() UserQuerySet {
	return qs.w(qs.db.Order("updated_at DESC"))
}

// SetCreatedAt is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetCreatedAt(createdAt time.Time) UserUpdater {
	u.fields[string(UserDBSchema.CreatedAt)] = createdAt
	return u
}

// SetFirstName is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetFirstName(firstName string) UserUpdater {
	u.fields[string(UserDBSchema.FirstName)] = firstName
	return u
}

// SetID is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetID(ID uint) UserUpdater {
	u.fields[string(UserDBSchema.ID)] = ID
	return u
}

// SetLastName is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetLastName(lastName string) UserUpdater {
	u.fields[string(UserDBSchema.LastName)] = lastName
	return u
}

// SetTelegramID is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetTelegramID(telegramID int64) UserUpdater {
	u.fields[string(UserDBSchema.TelegramID)] = telegramID
	return u
}

// SetUpdatedAt is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetUpdatedAt(updatedAt time.Time) UserUpdater {
	u.fields[string(UserDBSchema.UpdatedAt)] = updatedAt
	return u
}

// SetUsername is an autogenerated method
// nolint: dupl
func (u UserUpdater) SetUsername(username string) UserUpdater {
	u.fields[string(UserDBSchema.Username)] = username
	return u
}

// TelegramIDEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDEq(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id = ?", telegramID))
}

// TelegramIDGt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDGt(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id > ?", telegramID))
}

// TelegramIDGte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDGte(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id >= ?", telegramID))
}

// TelegramIDIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDIn(telegramID int64, telegramIDRest ...int64) UserQuerySet {
	iArgs := []interface{}{telegramID}
	for _, arg := range telegramIDRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("telegram_id IN (?)", iArgs))
}

// TelegramIDLt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDLt(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id < ?", telegramID))
}

// TelegramIDLte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDLte(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id <= ?", telegramID))
}

// TelegramIDNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDNe(telegramID int64) UserQuerySet {
	return qs.w(qs.db.Where("telegram_id != ?", telegramID))
}

// TelegramIDNotIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) TelegramIDNotIn(telegramID int64, telegramIDRest ...int64) UserQuerySet {
	iArgs := []interface{}{telegramID}
	for _, arg := range telegramIDRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("telegram_id NOT IN (?)", iArgs))
}

// Update is an autogenerated method
// nolint: dupl
func (u UserUpdater) Update() error {
	return u.db.Updates(u.fields).Error
}

// UpdatedAtEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtEq(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at = ?", updatedAt))
}

// UpdatedAtGt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtGt(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at > ?", updatedAt))
}

// UpdatedAtGte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtGte(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at >= ?", updatedAt))
}

// UpdatedAtLt is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtLt(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at < ?", updatedAt))
}

// UpdatedAtLte is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtLte(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at <= ?", updatedAt))
}

// UpdatedAtNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UpdatedAtNe(updatedAt time.Time) UserQuerySet {
	return qs.w(qs.db.Where("updated_at != ?", updatedAt))
}

// UsernameEq is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UsernameEq(username string) UserQuerySet {
	return qs.w(qs.db.Where("username = ?", username))
}

// UsernameIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UsernameIn(username string, usernameRest ...string) UserQuerySet {
	iArgs := []interface{}{username}
	for _, arg := range usernameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("username IN (?)", iArgs))
}

// UsernameNe is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UsernameNe(username string) UserQuerySet {
	return qs.w(qs.db.Where("username != ?", username))
}

// UsernameNotIn is an autogenerated method
// nolint: dupl
func (qs UserQuerySet) UsernameNotIn(username string, usernameRest ...string) UserQuerySet {
	iArgs := []interface{}{username}
	for _, arg := range usernameRest {
		iArgs = append(iArgs, arg)
	}
	return qs.w(qs.db.Where("username NOT IN (?)", iArgs))
}

// ===== END of query set UserQuerySet

// ===== BEGIN of User modifiers

type userDBSchemaField string

// UserDBSchema stores db field names of User
var UserDBSchema = struct {
	ID         userDBSchemaField
	CreatedAt  userDBSchemaField
	UpdatedAt  userDBSchemaField
	DeletedAt  userDBSchemaField
	TelegramID userDBSchemaField
	FirstName  userDBSchemaField
	LastName   userDBSchemaField
	Username   userDBSchemaField
}{

	ID:         userDBSchemaField("id"),
	CreatedAt:  userDBSchemaField("created_at"),
	UpdatedAt:  userDBSchemaField("updated_at"),
	DeletedAt:  userDBSchemaField("deleted_at"),
	TelegramID: userDBSchemaField("telegram_id"),
	FirstName:  userDBSchemaField("first_name"),
	LastName:   userDBSchemaField("last_name"),
	Username:   userDBSchemaField("username"),
}

// Update updates User fields by primary key
func (o *User) Update(db *gorm.DB, fields ...userDBSchemaField) error {
	dbNameToFieldName := map[string]interface{}{
		"id":          o.ID,
		"created_at":  o.CreatedAt,
		"updated_at":  o.UpdatedAt,
		"deleted_at":  o.DeletedAt,
		"telegram_id": o.TelegramID,
		"first_name":  o.FirstName,
		"last_name":   o.LastName,
		"username":    o.Username,
	}
	u := map[string]interface{}{}
	for _, f := range fields {
		fs := string(f)
		u[fs] = dbNameToFieldName[fs]
	}
	if err := db.Model(o).Updates(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}

		return fmt.Errorf("can't update User %v fields %v: %s",
			o, fields, err)
	}

	return nil
}

// UserUpdater is an User updates manager
type UserUpdater struct {
	fields map[string]interface{}
	db     *gorm.DB
}

// NewUserUpdater creates new User updater
func NewUserUpdater(db *gorm.DB) UserUpdater {
	return UserUpdater{
		fields: map[string]interface{}{},
		db:     db.Model(&User{}),
	}
}

// ===== END of User modifiers

// ===== END of all query sets
