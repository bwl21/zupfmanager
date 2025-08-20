 mchte # ABC File Directory Default Implementation

**Date:** 2025-08-19 19:40:49  
**Feature:** Default --abc-file-dir from most recent import  
**Branch:** feature/abc-file-dir-default  
**Commit:** 554ff12  

## Problem Statement

The user requested a default for the `--abc-file-dir` parameter in the project build command. This default should be set to the folder of the most recent import operation, providing a convenient workflow where users don't need to manually specify the ABC file directory after importing files.

## Solution Overview

Implemented a state tracking mechanism that:
1. Records the directory path whenever an import operation is performed
2. Uses this recorded path as the default for `--abc-file-dir` when not explicitly specified
3. Maintains backward compatibility with existing workflows

## Implementation Details

### Files Modified

#### 1. `pkg/core/import.go`

**Added constants:**
```go
const lastImportDirFile = ".zupfmanager_last_import_dir"
```

**Modified `ImportDirectory()` function:**
- Added call to `saveLastImportDir()` after successful import
- Includes error handling that logs warnings but doesn't fail the import

**New utility functions:**
```go
// saveLastImportDir saves the directory path to a state file
func saveLastImportDir(directory string) error

// GetLastImportDir retrieves the most recent import directory  
func GetLastImportDir() (string, error)
```

**Key features:**
- Stores absolute path in user's home directory (`~/.zupfmanager_last_import_dir`)
- Validates that the directory still exists when retrieving
- Returns empty string if no previous import or directory no longer exists

#### 2. `cmd/project-build.go`

**Modified default value logic:**
```go
if projectBuildAbcFileDir == "" {
    // Check if abc_file_dir exists and is a string
    abcFileDir, ok := project.Config["abc_file_dir"].(string)
    if ok {
        projectBuildAbcFileDir = abcFileDir
    } else {
        // Try to use the most recent import directory as default
        lastImportDir, err := core.GetLastImportDir()
        if err == nil && lastImportDir != "" {
            projectBuildAbcFileDir = lastImportDir
        } else {
            // Provide a default value or handle the error appropriately
            projectBuildAbcFileDir = ""
        }
    }
}
```

## Fallback Hierarchy

The system now uses the following priority order for determining the ABC file directory:

1. **Explicit flag:** `--abc-file-dir` parameter
2. **Project config:** `abc_file_dir` setting in project configuration
3. **Last import directory:** Most recent import directory from state file
4. **Empty string:** Original fallback behavior

## Testing Results

### Functional Testing

**Import tracking test:**
```bash
# Created test directory and ABC file
mkdir -p /tmp/test_abc_dir
echo 'T:Test Song...' > /tmp/test_abc_dir/test.abc

# Ran import command
/tmp/zupfmanager import /tmp/test_abc_dir

# Verified state file creation
cat ~/.zupfmanager_last_import_dir
# Output: /tmp/test_abc_dir
```

**Default retrieval test:**
```go
// Test program confirmed GetLastImportDir() returns correct path
lastDir, err := core.GetLastImportDir()
// Result: "/tmp/test_abc_dir"
```

**Build compilation test:**
```bash
go build -o /tmp/zupfmanager ./main.go
# Successful compilation with no errors
```

### Edge Cases Handled

1. **No previous imports:** Returns empty string, maintains original behavior
2. **Directory no longer exists:** Validates existence and returns empty string if missing
3. **File system errors:** Proper error handling with fallback to original behavior
4. **Concurrent access:** Uses atomic file operations for state persistence

## Backward Compatibility

- **Existing workflows:** No changes to current behavior when `--abc-file-dir` is explicitly specified
- **Project configurations:** Existing `abc_file_dir` settings in project configs take precedence
- **Error handling:** Graceful degradation when state file operations fail

## Benefits

1. **Improved UX:** Reduces need to repeatedly specify `--abc-file-dir` parameter
2. **Workflow efficiency:** Natural progression from import to build operations
3. **Zero configuration:** Works automatically without user setup
4. **Non-intrusive:** Doesn't affect existing established workflows

## Technical Considerations

### State File Location
- **Path:** `~/.zupfmanager_last_import_dir`
- **Format:** Plain text file containing absolute directory path
- **Permissions:** 0644 (readable by user and group, writable by user)

### Error Handling Strategy
- **Import failures:** Log warnings but don't fail import operations
- **State read failures:** Silently fall back to original behavior
- **Directory validation:** Check existence before using as default

### Performance Impact
- **Minimal overhead:** Single file read/write operations
- **Lazy evaluation:** State only read when needed for build operations
- **No database changes:** Uses simple file-based state management

## Future Enhancements

Potential improvements that could be considered:
1. **Multiple directory tracking:** Remember last N import directories
2. **Project-specific defaults:** Per-project import directory preferences
3. **Configuration option:** Allow users to disable this feature
4. **Import history:** Full audit trail of import operations

## Git Information

**Branch:** `feature/abc-file-dir-default`  
**Commit:** `554ff12`  
**Commit Message:**
```
feat: add default abc-file-dir from most recent import

- Track import directory in ~/.zupfmanager_last_import_dir state file
- Use most recent import directory as default for --abc-file-dir when not specified
- Fallback hierarchy: explicit flag > project config > last import dir > empty
- Backward compatible with existing workflows

Co-authored-by: Ona <no-reply@ona.com>
```

**Files Changed:**
- `cmd/project-build.go` (6 lines added, 2 lines modified)
- `pkg/core/import.go` (49 lines added)

## Directory Picker Widget Enhancement

**Update:** 2025-08-20 10:23  
**Additional Commit:** `f1026a1`

## Database Schema Separation

**Update:** 2025-08-20 11:04  
**Additional Commit:** `fb08ad3`

### Critical Architecture Fix

The initial implementation incorrectly stored the `abc_file_dir_preference` in the project's `config` field. However, this field is specifically designed for Zupfnoter configuration templates, not zupfmanager internal settings.

#### Schema Changes:
1. **New Database Field:** Added `abc_file_dir_preference` as dedicated string field in project schema
2. **Clean Separation:** Project `config` now only contains Zupfnoter-specific settings
3. **Database Migration:** Required for existing installations

#### Technical Implementation:
- **Schema:** `internal/ent/schema/project.go` - Added dedicated field
- **API Models:** Updated `ProjectResponse` to include new field
- **Converters:** Updated domain model conversion functions
- **Handlers:** Modified all project API endpoints to handle new field

#### Benefits:
- **Clean Architecture:** Clear separation between zupfmanager and Zupfnoter settings
- **Data Integrity:** Zupfnoter config remains untouched by UI preferences
- **Maintainability:** Easier to manage different types of configuration
- **Future-Proof:** Allows for additional zupfmanager-specific settings

## Full Path Requirement Fix

**Update:** 2025-08-20 11:35  
**Additional Commit:** `0bf2104`

### Critical Path Handling Issue

The initial directory picker implementation only stored directory names (e.g., `music-folder`) instead of full absolute paths (e.g., `/home/user/music-folder`). This caused failures in Zupfnoter, which requires complete file system paths.

#### Root Cause:
Browser security restrictions prevent direct access to full file system paths. The `webkitRelativePath` API only provides relative paths within the selected directory structure.

#### Solution Implemented:
1. **Hybrid Approach:** Directory picker helps locate directories, but users must enter full paths manually
2. **Path Validation:** Visual warnings for incomplete or invalid paths
3. **Multi-Platform Support:** Handles Unix/Linux (`/path`) and Windows (`C:\path`) formats
4. **Smart Saving:** Only saves preferences when valid absolute paths are provided

#### UX Improvements:
- **Clear Instructions:** Placeholder text shows expected path format
- **Visual Feedback:** Directory selection shows found ABC files count
- **Path Validation:** Real-time warnings for incomplete paths
- **Examples:** Shows platform-specific path examples in warnings

#### Technical Implementation:
- **Path Validation Function:** `isValidPath()` checks for absolute path patterns
- **Directory Info Display:** Shows selected directory name and ABC file count
- **Conditional Saving:** Only saves valid paths as preferences
- **Cross-Platform:** Supports Windows, macOS, and Linux path formats

This ensures Zupfnoter receives the complete paths it needs for successful ABC file processing.

### Enhanced Frontend Implementation

The solution was further enhanced with a proper directory picker widget and database-based persistence:

#### Frontend Improvements:
1. **Directory Picker Widget:**
   - Browse button with folder icon
   - Support for modern File System Access API
   - Fallback to `webkitdirectory` for broader browser support
   - Visual loading states and improved UX

2. **Database Persistence:**
   - Removed localStorage dependency
   - Directory preference stored in project configuration
   - Multi-device consistency
   - Automatic saving when directory is selected

#### Backend Enhancements:
1. **New API Endpoint:** `PUT /api/v1/projects/{id}/abc-file-dir`
   - Updates `abc_file_dir_preference` in project config
   - Returns updated project configuration

2. **Enhanced Priority System:**
   ```
   1. Explicit --abc-file-dir flag
   2. Project abc_file_dir_preference (user selection)
   3. Project abc_file_dir config
   4. Last import directory
   5. Empty string
   ```

#### Technical Implementation:
- **Frontend:** Modern directory selection with File System Access API fallback
- **Backend:** Project config field `abc_file_dir_preference`
- **Persistence:** Database-based storage in project configuration JSON
- **API:** RESTful endpoint for preference updates

#### User Experience:
1. User clicks "Browse" button in build modal
2. System opens directory picker (native or web-based)
3. Selected directory is immediately saved to project configuration
4. Directory appears as default in future build operations
5. Preference persists across sessions and devices

## Conclusion

The implementation successfully provides a convenient default for the `--abc-file-dir` parameter while maintaining full backward compatibility. The enhanced solution includes:

- **Intuitive UI:** Directory picker widget with browse functionality
- **Robust Persistence:** Database-based storage for multi-device consistency  
- **Modern Browser Support:** File System Access API with graceful fallbacks
- **Flexible Priority System:** Multiple fallback levels for different use cases
- **Seamless Integration:** Works with existing project configuration system

The solution is robust, handles edge cases appropriately, and integrates seamlessly with the existing codebase architecture while providing a significantly improved user experience.
