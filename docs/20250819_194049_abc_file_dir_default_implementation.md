 mchte  # ABC File Directory Default Implementation

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

## Browse Button Removal

**Update:** 2025-08-20 18:00  
**Additional Commits:** `920859f`, `eae5af2`, `8c0193f`, `8c284c0`, `0ca79fd`

### Critical Browser Security Limitation

After implementing directory picker widgets across the application (ImportView, BuildConfigModal, ProjectsView), a fundamental limitation was discovered: **Browser security restrictions prevent access to absolute file system paths**, which are required by the backend for Zupfnoter integration.

#### Problem Identified:
1. **Browser APIs Limitation:** `webkitRelativePath` and File System Access API only provide relative paths or directory names
2. **Backend Requirement:** Zupfnoter requires absolute paths (e.g., `/home/user/Documents/music/`)
3. **User Confusion:** Browse buttons suggested functionality that couldn't deliver absolute paths
4. **Inconsistent UX:** Buttons that appeared functional but required manual path completion

#### Root Cause Analysis:
```javascript
// Browser APIs can only provide:
file.webkitRelativePath  // "music-folder/song.abc" (relative)
dirHandle.name          // "music-folder" (name only)

// But backend needs:
"/home/user/Documents/music-folder/"  // Absolute path - NOT accessible via browser
```

#### Solution: Complete Browse Button Removal

**Decision:** Remove all browse buttons and directory picker functionality from the entire application to provide honest, consistent UX.

#### Files Modified:
1. **frontend/src/views/ImportView.vue**
   - Removed file and directory browse buttons
   - Removed hidden input elements and File APIs
   - Simplified to text-only path input
   - Removed all browser-based file selection functions

2. **frontend/src/components/BuildConfigModal.vue**
   - Removed directory browse button for ABC file directory
   - Removed File System Access API integration
   - Cleaned up directory selection display components
   - Simplified to manual path entry

3. **frontend/src/views/ProjectsView.vue**
   - Removed directory browse button for project creation
   - Removed directory selection validation
   - Cleaned up related UI components and functions

#### Technical Changes:
- **Removed Functions:** `openDirectoryPicker()`, `handleDirectorySelection()`, `selectFile()`, `selectDirectory()`
- **Removed Refs:** `fileInput`, `directoryInput`, `selectedFileInfo`, `selectedDirectoryInfo`
- **Removed UI Elements:** Browse buttons, file selection dialogs, directory info displays
- **Cleaned Text:** Removed all references to "Browse" functionality in help text and comments

#### Commit History:
1. **920859f** - `feat: add file and directory browser to import view`
   - Initial implementation of browse buttons with File System Access API
   - Added file/directory selection dialogs and validation
   
2. **eae5af2** - `fix: improve file/directory browser path handling`
   - Attempted to fix path display issues
   - Added better user guidance and validation
   
3. **8c0193f** - `revert: remove file/directory browser functionality`
   - Removed browse functionality from ImportView
   - Recognized browser security limitations
   
4. **8c284c0** - `remove: eliminate all browse buttons from application`
   - Systematic removal from BuildConfigModal and ProjectsView
   - Complete cleanup of all browser-based file selection
   
5. **0ca79fd** - `cleanup: remove remaining browse button references`
   - Final cleanup of help text and comments
   - Complete removal of all browse-related UI text

#### Bundle Size Impact:
- **ImportView:** Reduced from 12+ kB to 8.63 kB
- **ProjectsView:** Reduced from 11+ kB to 9.55 kB  
- **BuildConfigModal:** Reduced bundle size through removed File API code

#### User Experience Changes:
- **Before:** Browse buttons that couldn't provide absolute paths → user confusion
- **After:** Clear text input fields with absolute path requirements → honest UX
- **Benefit:** No false promises of functionality that browser security prevents

#### Alternative Solutions Considered:
1. **Server-side file browser:** Would require backend file system access
2. **Relative path resolution:** Would need configurable base directories
3. **Path completion hints:** Still wouldn't solve absolute path requirement

#### Final Architecture:
```
User Input: Manual absolute path entry
├── ImportView: /path/to/abc/files/
├── BuildConfigModal: /path/to/abc/directory/  
└── ProjectsView: /path/to/project/directory/
```

## Conclusion

The implementation successfully provides a convenient default for the `--abc-file-dir` parameter while maintaining full backward compatibility. The final solution includes:

- **Honest UX:** Text-only path input without misleading browse buttons
- **Robust Persistence:** Database-based storage for multi-device consistency  
- **Clean Architecture:** Removed unused browser APIs and file selection code
- **Flexible Priority System:** Multiple fallback levels for different use cases
- **Seamless Integration:** Works with existing project configuration system
- **Reduced Bundle Size:** Smaller JavaScript bundles through code removal

**Key Lesson:** Browser security restrictions make absolute path access impossible, requiring manual path entry for backend systems that need complete file system paths. The solution prioritizes honest user experience over feature appearance.

## PDF Preview Functionality

**Update:** 2025-08-20 18:30  
**Additional Commits:** TBD

### Preview System Implementation

Added comprehensive PDF preview functionality that allows users to generate and view PDF outputs for individual ABC files without running a full project build.

#### Architecture Overview:
```
Frontend (Vue.js) → API Endpoints → Song Service → Zupfnoter → PDF Generation
                                                  ↓
                                            Temporary Storage → PDF Serving
```

#### Backend Implementation:

**1. Extended Song Service Interface:**
```go
type SongService interface {
    // Existing methods...
    
    // Preview operations
    GeneratePreview(ctx context.Context, req GeneratePreviewRequest) (*GeneratePreviewResponse, error)
    ListPreviewPDFs(ctx context.Context, songID int) ([]*PreviewPDF, error)
    GetPreviewPDF(ctx context.Context, songID int, filename string) (string, error)
    CleanupPreviewPDFs(ctx context.Context, songID int) error
}
```

**2. New API Endpoints:**
- `POST /api/v1/songs/{id}/generate-preview` - Generate preview PDFs
- `GET /api/v1/songs/{id}/preview-pdfs` - List available PDFs
- `GET /api/v1/songs/{id}/preview-pdf/{filename}` - Serve PDF file
- `DELETE /api/v1/songs/{id}/preview-pdfs` - Cleanup preview PDFs

**3. Preview Storage Strategy:**
- **Location:** `/tmp/zupfmanager/previews/song-{id}/pdf/`
- **Naming:** `{abc_filename}*.pdf` (multiple variants per song)
- **Lifecycle:** User-controlled cleanup + temporary storage
- **Security:** Path traversal protection, filename validation

#### Frontend Implementation:

**1. Preview Modal Component:**
```vue
<PreviewModal>
  <!-- Generate Section -->
  <input v-model="abcFileDir" placeholder="/path/to/abc/files" />
  <button @click="generatePreview">Generate Preview PDFs</button>
  
  <!-- PDF List -->
  <div v-for="pdf in pdfs" class="pdf-item">
    <span>{{ pdf.filename }}</span>
    <button @click="openPDF(pdf.filename)">Open</button>
  </div>
</PreviewModal>
```

**2. Integration with SongDetailView:**
- Added "Preview PDFs" button in actions section
- Modal-based interface for preview management
- Real-time PDF list updates after generation

**3. API Service Extensions:**
```typescript
export const songApi = {
  // Existing methods...
  
  generatePreview: (id: number, data: { abc_file_dir: string }) => Promise<GeneratePreviewResponse>
  listPreviewPDFs: (id: number) => Promise<PreviewPDFListResponse>
  getPreviewPDFUrl: (id: number, filename: string) => string
  cleanupPreviewPDFs: (id: number) => Promise<MessageResponse>
}
```

#### Technical Features:

**1. PDF Generation Process:**
1. User selects song and provides ABC file directory
2. System creates temporary preview directory
3. Zupfnoter generates multiple PDF variants (A3, M, O, B, X layouts)
4. PDFs stored in organized temporary structure
5. File list returned to frontend for display

**2. PDF Serving:**
- Direct file serving with proper MIME types
- Security validation (no path traversal)
- Opens in new browser tab for viewing
- Supports all Zupfnoter-generated PDF variants

**3. Lifecycle Management:**
- User-initiated cleanup via "Clear All" button
- Temporary storage in system temp directory
- Automatic cleanup on system restart
- Per-song isolation prevents conflicts

#### User Experience:

**1. Preview Workflow:**
1. Navigate to song detail page
2. Click "Preview PDFs" button
3. Enter ABC file directory path
4. Click "Generate Preview PDFs"
5. View generated PDFs in list
6. Click "Open" to view PDF in new tab

**2. Management Features:**
- **Refresh:** Update PDF list without regeneration
- **Clear All:** Remove all preview PDFs for song
- **Real-time Updates:** List updates after generation/cleanup
- **Error Handling:** Clear error messages for failed operations

#### Benefits:

1. **Quick Preview:** View PDF outputs without full project build
2. **Multiple Variants:** See all Zupfnoter layout options
3. **No Project Required:** Works with individual ABC files
4. **Temporary Storage:** No permanent disk usage
5. **User Control:** Manual cleanup and regeneration

#### Technical Considerations:

**1. Performance:**
- Isolated PDF generation per song
- Temporary storage prevents accumulation
- Concurrent generation support via goroutines

**2. Security:**
- Path traversal protection in filename handling
- Temporary directory isolation
- No permanent file system modifications

**3. Error Handling:**
- Zupfnoter execution error capture
- File system error handling
- User-friendly error messages in UI

The solution is robust, handles edge cases appropriately, and integrates seamlessly with the existing codebase architecture while providing a clear, unambiguous user experience that works within browser security constraints.
