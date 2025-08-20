# User Guide: Preview Functionality

**Last Updated:** 2025-08-20  
**Version:** 1.0  

## Overview

The preview functionality allows you to quickly view PDF outputs for individual songs in your projects. The system automatically uses your project's ABC file directory setting to find and display available PDFs.

## How to Use Preview

### 1. Individual Song Preview

Each song in your project now has its own **Preview** button:

1. Navigate to a project (Projects → Select Project)
2. In the "Songs in Project" section, find the song you want to preview
3. Click the green **Preview** button next to the song
4. The preview modal will open automatically

### 2. Automatic Directory Detection

When you open a preview:

- The system automatically uses your project's **ABC File Directory** setting
- If PDFs exist in that directory, they will be found automatically
- No manual directory entry required in most cases

### 3. Manual Directory Override

If you need to use a different directory:

1. In the preview modal, modify the "ABC File Directory" field
2. Click **Find PDFs** to search the new location
3. The system will remember this override for the current session

## Button Guide

In the song list, you'll see three action buttons for each song:

- 🟢 **Preview** - View available PDFs for this song
- 🔵 **Edit** - Modify song settings (difficulty, priority, comments)
- 🔴 **Remove** - Remove song from project

## PDF Discovery

The preview system looks for PDFs using this pattern:
- **Location:** Your project's ABC file directory
- **Naming:** Files that start with the song's ABC filename
- **Examples:** 
  - `song.abc` → finds `song.pdf`, `song_A4.pdf`, `song_layout2.pdf`
  - `melody.abc` → finds `melody.pdf`, `melody_compact.pdf`

## Project Settings Integration

### Setting ABC File Directory

1. Go to Projects → Select Project → **Edit Project**
2. Set the **ABC File Directory Preference**
3. Save the project
4. All song previews will now use this directory automatically

### Benefits of Setting Directory Preference

- **Automatic Preview:** No manual directory entry needed
- **Consistency:** All songs in the project use the same directory
- **Time Saving:** Faster access to previews
- **Error Reduction:** Less chance of typos in directory paths

## Troubleshooting

### No PDFs Found

**Problem:** Preview shows "No PDFs found"

**Solutions:**
1. Verify the ABC file directory path is correct
2. Check that PDF files exist in that directory
3. Ensure PDF filenames start with the ABC filename
4. Try clicking **Refresh PDFs** to re-scan

### Wrong Directory

**Problem:** Preview is looking in the wrong directory

**Solutions:**
1. Update the project's ABC File Directory Preference
2. Or manually enter the correct path in the preview modal
3. The manual path will be used for the current session

### Permission Errors

**Problem:** "Cannot access directory" error

**Solutions:**
1. Check that the directory path exists
2. Verify you have read permissions for the directory
3. Ensure the path is absolute (starts with `/` on Unix or `C:\` on Windows)

## Workflow Tips

### Efficient Preview Workflow

1. **Set Project Directory:** Configure ABC File Directory Preference once per project
2. **Direct Preview:** Click preview buttons directly from song list
3. **Multiple Variants:** View all PDF variants created by external tools
4. **Quick Access:** No need to navigate away from the song list

### Integration with External Tools

The preview system works well with:
- **Zupfnoter:** Interactive ABC editor and PDF generator
- **abc2ps/abcm2ps:** Command-line ABC to PDF converters
- **Custom Scripts:** Any tool that creates PDFs from ABC files

Just ensure your external tools save PDFs in the same directory as your ABC files with matching filenames.

## Best Practices

1. **Consistent Directory Structure:** Keep ABC files and PDFs in the same directory
2. **Clear Naming:** Use descriptive ABC filenames that match your song titles
3. **Project Settings:** Set ABC File Directory Preference for each project
4. **Regular Cleanup:** Remove old PDF variants you no longer need

## Advanced Features

### Multiple PDF Variants

The system can display multiple PDF versions of the same song:
- `song.pdf` - Default layout
- `song_A4.pdf` - A4 paper size
- `song_compact.pdf` - Compact layout
- `song_large.pdf` - Large print version

All variants will appear in the preview list, allowing you to choose the best version for your needs.

### Refresh Functionality

Use the **Refresh PDFs** button to:
- Re-scan the directory after creating new PDFs
- Update the list after external tool processing
- Refresh after moving or renaming PDF files

This ensures you always see the most current PDF files available.
