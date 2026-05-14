// Package redactor provides utilities for masking sensitive values
// in environment variable maps before display or export.
//
// Keys matching patterns such as PASSWORD, SECRET, TOKEN, or API_KEY
// have their values replaced with [REDACTED]. Custom patterns can be
// supplied when creating a Redactor via New.
package redactor
