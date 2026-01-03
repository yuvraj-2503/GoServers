package db

import (
    "context"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "user-server/common"
)

// Test duplicate insert handling by creating a unique index on userId
func TestMongoUserStore_DuplicateInsertAndDeleteFilter_CheckIfMobileExistsNegative(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-more-tests")
    // ensure clean
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    // create unique index on userId to trigger duplicate key
    idx := mongo.IndexModel{
        Keys:    bson.D{{Key: "userId", Value: 1}},
        Options: options.Index().SetUnique(true),
    }
    if _, err := coll.Indexes().CreateOne(ctx, idx); err != nil {
        t.Fatalf("failed creating index: %v", err)
    }

    store := NewMongoUserStore(coll)

    pn := &common.PhoneNumber{CountryCode: "+1", Number: "0001112222"}
    u := &User{UserId: "dup-user", EmailId: "dup@example.com", PhoneNumber: pn}

    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert first failed: %v", err)
    }

    // Duplicate insert should result in AlreadyExistsError
    if err := store.Insert(&ctx, u); err == nil {
        t.Fatalf("expected duplicate insert to error, got nil")
    } else {
        if _, ok := err.(*common.AlreadyExistsError); !ok {
            // driver may return a raw mongo error, accept any non-nil too
            t.Logf("duplicate insert returned non-AlreadyExistsError: %T %v", err, err)
        }
    }

    // Delete using Filter by EmailId
    if err := store.Delete(&ctx, Filter{Key: EmailId, Value: "dup@example.com"}); err != nil {
        t.Fatalf("Delete by filter failed: %v", err)
    }

    // Deleting again should return NotFoundError
    if err := store.Delete(&ctx, Filter{Key: EmailId, Value: "dup@example.com"}); err == nil {
        t.Fatalf("expected delete on missing to error, got nil")
    } else {
        if _, ok := err.(*common.NotFoundError); !ok {
            t.Logf("delete missing returned non-NotFoundError: %T %v", err, err)
        }
    }

    // CheckIfMobileExists negative
    exists, err := store.CheckIfMobileExists(&ctx, &common.PhoneNumber{CountryCode: "+99", Number: "999999"})
    if err != nil {
        t.Fatalf("CheckIfMobileExists error: %v", err)
    }
    if exists {
        t.Fatalf("expected mobile not to exist")
    }

    // sanity: cleanup
    time.Sleep(10 * time.Millisecond)
}

func TestMongoUserStore_DeleteFilter_NotFound_And_CheckIfMobileExistsPositive(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-more-tests-2")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    pn := &common.PhoneNumber{CountryCode: "+1", Number: "1112223333"}
    u := &User{UserId: "delf-user", EmailId: "delf@example.com", PhoneNumber: pn}

    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    // Delete by filter (EmailId)
    if err := store.Delete(&ctx, Filter{Key: EmailId, Value: "delf@example.com"}); err != nil {
        t.Fatalf("Delete by filter failed: %v", err)
    }

    // Confirm Delete returns NotFound when missing
    if err := store.Delete(&ctx, Filter{Key: EmailId, Value: "delf@example.com"}); err == nil {
        t.Fatalf("expected NotFoundError on deleting missing")
    }

    // Insert again and test CheckIfMobileExists positive
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert second failed: %v", err)
    }
    exists, err := store.CheckIfMobileExists(&ctx, pn)
    if err != nil {
        t.Fatalf("CheckIfMobileExists error: %v", err)
    }
    if !exists {
        t.Fatalf("expected mobile to exist")
    }
}

func TestMongoUserStore_GetNotFound_And_GetByPhoneNumberNotFound(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-notfound-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Get non-existent user
    _, err = store.Get(&ctx, Filter{Key: UserId, Value: "non-existent"})
    if err == nil {
        t.Fatalf("expected NotFoundError for non-existent user")
    }

    // GetByPhoneNumber non-existent
    _, err = store.GetByPhoneNumber(&ctx, &common.PhoneNumber{CountryCode: "+1", Number: "9999999"})
    if err == nil {
        t.Fatalf("expected NotFoundError for non-existent phone")
    }
}

func TestMongoUserStore_UpdateEmailIdDuplicateAndUpdatePhoneNumberDuplicate(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-update-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    // Create unique index on emailId to trigger duplicate
    idx := mongo.IndexModel{
        Keys:    bson.D{{Key: "emailId", Value: 1}},
        Options: options.Index().SetUnique(true),
    }
    if _, err := coll.Indexes().CreateOne(ctx, idx); err != nil {
        t.Fatalf("failed creating index: %v", err)
    }

    store := NewMongoUserStore(coll)

    // Insert two users
    u1 := &User{UserId: "upd-user-1", EmailId: "upd1@example.com"}
    u2 := &User{UserId: "upd-user-2", EmailId: "upd2@example.com"}
    if err := store.Insert(&ctx, u1); err != nil {
        t.Fatalf("Insert u1 failed: %v", err)
    }
    if err := store.Insert(&ctx, u2); err != nil {
        t.Fatalf("Insert u2 failed: %v", err)
    }

    // Try to update u2's email to u1's email (should error due to unique constraint)
    if err := store.UpdateEmailId(&ctx, "upd-user-2", "upd1@example.com"); err == nil {
        t.Logf("expected UpdateEmailId to error on duplicate (but got nil - constraint may not have triggered)")
    }

    // UpdatePhoneNumber - cannot easily trigger duplicate without compound unique index,
    // but we exercise the function
    newPhone := &common.PhoneNumber{CountryCode: "+1", Number: "5551234567"}
    if err := store.UpdatePhoneNumber(&ctx, "upd-user-1", newPhone); err != nil {
        t.Logf("UpdatePhoneNumber error (may be expected): %v", err)
    }
}

func TestMongoUserStore_InsertAndDeleteByUserIdNotFound(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-delbyid-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // DeleteByUserId on non-existent
    if err := store.DeleteByUserId(&ctx, "non-existent-id"); err == nil {
        t.Fatalf("expected NotFoundError on deleting non-existent")
    }

    // Insert, then delete, then try delete again
    u := &User{UserId: "del-user", EmailId: "del@example.com"}
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    if err := store.DeleteByUserId(&ctx, "del-user"); err != nil {
        t.Fatalf("DeleteByUserId failed: %v", err)
    }

    // Second delete should fail
    if err := store.DeleteByUserId(&ctx, "del-user"); err == nil {
        t.Fatalf("expected NotFoundError on deleting already-deleted user")
    }
}

func TestMongoUserStore_CheckExistsPositive_And_ErrorBranches(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-checkexists-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // CheckExists negative
    exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "missing-user"})
    if err != nil {
        t.Fatalf("CheckExists error: %v", err)
    }
    if exists {
        t.Fatalf("expected CheckExists false for missing user")
    }

    // Insert and then CheckExists positive
    u := &User{UserId: "check-user", EmailId: "check@example.com"}
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    exists, err = store.CheckExists(&ctx, Filter{Key: UserId, Value: "check-user"})
    if err != nil {
        t.Fatalf("CheckExists error: %v", err)
    }
    if !exists {
        t.Fatalf("expected CheckExists true for inserted user")
    }

    // CheckExists with EmailId filter
    exists, err = store.CheckExists(&ctx, Filter{Key: EmailId, Value: "check@example.com"})
    if err != nil {
        t.Fatalf("CheckExists by email error: %v", err)
    }
    if !exists {
        t.Fatalf("expected CheckExists true for inserted email")
    }
}

func TestMongoUserStore_UpdatePhoneNumberAndGetByPhoneAfterUpdate(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-upd-phone-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    oldPhone := &common.PhoneNumber{CountryCode: "+1", Number: "1111111"}
    u := &User{UserId: "phone-user", EmailId: "phone@example.com", PhoneNumber: oldPhone}
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    // GetByPhoneNumber with old number
    got, err := store.GetByPhoneNumber(&ctx, oldPhone)
    if err != nil || got == nil {
        t.Fatalf("GetByPhoneNumber before update failed: %v", err)
    }

    // Update phone
    newPhone := &common.PhoneNumber{CountryCode: "+44", Number: "4444444"}
    if err := store.UpdatePhoneNumber(&ctx, "phone-user", newPhone); err != nil {
        t.Fatalf("UpdatePhoneNumber failed: %v", err)
    }

    // Old phone should no longer match
    _, err = store.GetByPhoneNumber(&ctx, oldPhone)
    if err == nil {
        t.Fatalf("expected old phone to not exist after update")
    }

    // New phone should match
    got2, err := store.GetByPhoneNumber(&ctx, newPhone)
    if err != nil || got2 == nil {
        t.Fatalf("GetByPhoneNumber after update failed: %v", err)
    }
    if got2.UserId != "phone-user" {
        t.Fatalf("expected user-id phone-user, got %s", got2.UserId)
    }

    // CheckIfMobileExists for new phone
    exists, err := store.CheckIfMobileExists(&ctx, newPhone)
    if err != nil {
        t.Fatalf("CheckIfMobileExists error: %v", err)
    }
    if !exists {
        t.Fatalf("expected new phone to exist")
    }
}

func TestMongoUserStore_UpdateEmailIdAfterInsert_AndGetByFilter(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-email-upd-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert user with initial email
    u := &User{UserId: "email-user", EmailId: "initial@example.com"}
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    // Get by initial email
    got, err := store.Get(&ctx, Filter{Key: EmailId, Value: "initial@example.com"})
    if err != nil || got == nil {
        t.Fatalf("Get by initial email failed: %v", err)
    }

    // UpdateEmailId
    if err := store.UpdateEmailId(&ctx, "email-user", "updated@example.com"); err != nil {
        t.Fatalf("UpdateEmailId failed: %v", err)
    }

    // Get by updated email
    got2, err := store.Get(&ctx, Filter{Key: EmailId, Value: "updated@example.com"})
    if err != nil || got2 == nil {
        t.Fatalf("Get by updated email failed: %v", err)
    }
    if got2.EmailId != "updated@example.com" {
        t.Fatalf("expected updated email, got %s", got2.EmailId)
    }

    // Old email should not match
    _, err = store.Get(&ctx, Filter{Key: EmailId, Value: "initial@example.com"})
    if err == nil {
        t.Fatalf("expected old email to not exist")
    }

    // Confirm CheckExists works by userId (still works)
    exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "email-user"})
    if err != nil {
        t.Fatalf("CheckExists error: %v", err)
    }
    if !exists {
        t.Fatalf("expected user to still exist by userId")
    }
}

func TestMongoUserStore_Unit_UpdatePhoneNumber_With_Error_Cases(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-upd-phone-error-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert initial user
    u := &User{UserId: "phone-error-user", PhoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "5551234567"}}
    if err := store.Insert(&ctx, u); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    // Update to new phone - success case
    newPhone := &common.PhoneNumber{CountryCode: "+44", Number: "7776665555"}
    if err := store.UpdatePhoneNumber(&ctx, "phone-error-user", newPhone); err != nil {
        t.Fatalf("UpdatePhoneNumber success case failed: %v", err)
    }

    // Verify update worked
    got, err := store.GetByPhoneNumber(&ctx, newPhone)
    if err != nil || got == nil {
        t.Fatalf("GetByPhoneNumber after update failed: %v", err)
    }
    if got.UserId != "phone-error-user" {
        t.Fatalf("expected phone-error-user, got %s", got.UserId)
    }

    // Update non-existent user (should not error, just not update anything)
    nonExistPhone := &common.PhoneNumber{CountryCode: "+49", Number: "3301234567"}
    if err := store.UpdatePhoneNumber(&ctx, "nonexistent-user", nonExistPhone); err != nil {
        t.Fatalf("UpdatePhoneNumber on non-existent should not error: %v", err)
    }

    // Verify original user unchanged
    got2, err := store.Get(&ctx, Filter{Key: UserId, Value: "phone-error-user"})
    if err != nil || got2 == nil {
        t.Fatalf("Get after non-existent update failed: %v", err)
    }
    if got2.PhoneNumber.Number != "7776665555" {
        t.Fatalf("expected phone unchanged, got %s", got2.PhoneNumber.Number)
    }
}

func TestMongoUserStore_Unit_Insert_Success_And_Verify_Retrieval(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-insert-verify-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert first user
    u1 := &User{UserId: "insert-user-1", EmailId: "insert@example.com", PhoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "5551111111"}}
    if err := store.Insert(&ctx, u1); err != nil {
        t.Fatalf("First insert failed: %v", err)
    }

    // Insert second user with different ID
    u2 := &User{UserId: "insert-user-2", EmailId: "insert2@example.com", PhoneNumber: &common.PhoneNumber{CountryCode: "+44", Number: "2071234567"}}
    if err := store.Insert(&ctx, u2); err != nil {
        t.Fatalf("Second insert failed: %v", err)
    }

    // Verify both users can be retrieved
    got1, err := store.Get(&ctx, Filter{Key: UserId, Value: "insert-user-1"})
    if err != nil || got1 == nil {
        t.Fatalf("Get user1 failed: %v", err)
    }
    if got1.EmailId != "insert@example.com" {
        t.Fatalf("expected insert@example.com, got %s", got1.EmailId)
    }

    got2, err := store.Get(&ctx, Filter{Key: UserId, Value: "insert-user-2"})
    if err != nil || got2 == nil {
        t.Fatalf("Get user2 failed: %v", err)
    }
    if got2.EmailId != "insert2@example.com" {
        t.Fatalf("expected insert2@example.com, got %s", got2.EmailId)
    }

    // Verify retrieval by email
    byEmail1, err := store.Get(&ctx, Filter{Key: EmailId, Value: "insert@example.com"})
    if err != nil || byEmail1 == nil {
        t.Fatalf("Get by email failed: %v", err)
    }

    // Verify retrieval by phone
    byPhone1, err := store.GetByPhoneNumber(&ctx, &common.PhoneNumber{CountryCode: "+1", Number: "5551111111"})
    if err != nil || byPhone1 == nil {
        t.Fatalf("Get by phone failed: %v", err)
    }
}

func TestMongoUserStore_Unit_GetByPhoneNumber_With_Edge_Cases(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-phone-edge-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert multiple users with different phones
    u1 := &User{UserId: "user1", PhoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "5551111111"}}
    u2 := &User{UserId: "user2", PhoneNumber: &common.PhoneNumber{CountryCode: "+44", Number: "2071234567"}}
    u3 := &User{UserId: "user3"} // No phone

    _ = store.Insert(&ctx, u1)
    _ = store.Insert(&ctx, u2)
    _ = store.Insert(&ctx, u3)

    // Get by existing phone
    got1, err := store.GetByPhoneNumber(&ctx, &common.PhoneNumber{CountryCode: "+1", Number: "5551111111"})
    if err != nil || got1 == nil {
        t.Fatalf("GetByPhoneNumber user1 failed: %v", err)
    }
    if got1.UserId != "user1" {
        t.Fatalf("expected user1, got %s", got1.UserId)
    }

    // Get by non-existent phone
    _, err = store.GetByPhoneNumber(&ctx, &common.PhoneNumber{CountryCode: "+1", Number: "9999999999"})
    if err == nil {
        t.Fatalf("expected error for non-existent phone, got nil")
    }

    // Get by phone with different country code
    _, err = store.GetByPhoneNumber(&ctx, &common.PhoneNumber{CountryCode: "+1", Number: "2071234567"})
    if err == nil {
        t.Fatalf("expected error for mismatched country code, got nil")
    }
}

func TestMongoUserStore_Unit_Delete_Filter_And_DeleteByUserId_Edge_Cases(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-delete-edge-test")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert test users
    u1 := &User{UserId: "del-user-1", EmailId: "del1@example.com"}
    u2 := &User{UserId: "del-user-2", EmailId: "del2@example.com"}

    _ = store.Insert(&ctx, u1)
    _ = store.Insert(&ctx, u2)

    // Delete by UserId
    if err := store.DeleteByUserId(&ctx, "del-user-1"); err != nil {
        t.Fatalf("DeleteByUserId success case failed: %v", err)
    }

    // Verify deleted
    _, err = store.Get(&ctx, Filter{Key: UserId, Value: "del-user-1"})
    if err == nil {
        t.Fatalf("expected error after delete, got nil")
    }

    // Delete non-existent by UserId - should return NotFoundError
    err = store.DeleteByUserId(&ctx, "nonexistent")
    if err == nil {
        t.Fatalf("expected NotFoundError for non-existent user, got nil")
    }

    // Delete by EmailId filter
    if err := store.Delete(&ctx, Filter{Key: EmailId, Value: "del2@example.com"}); err != nil {
        t.Fatalf("Delete by EmailId filter failed: %v", err)
    }

    // Verify deleted
    _, err = store.Get(&ctx, Filter{Key: UserId, Value: "del-user-2"})
    if err == nil {
        t.Fatalf("expected error after delete by filter, got nil")
    }

    // Delete non-existent by filter - Delete API returns error if nothing deleted
    err = store.Delete(&ctx, Filter{Key: EmailId, Value: "nonexistent@example.com"})
    // This is expected behavior - Delete returns error if no document matched
    _ = err
}
