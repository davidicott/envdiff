// Package merger provides functionality for combining two parsed .env maps
// into a single unified map.
//
// It supports three conflict-resolution strategies:
//
//   - PreferA: values from the first file win on conflict.
//   - PreferB: values from the second file win on conflict.
//   - ErrorOnConflict: returns an error if any key has differing values.
//
// Conflicts are always recorded in the Result regardless of strategy
// (except when ErrorOnConflict aborts early).
package merger
