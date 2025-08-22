# User Guide: Preview Functionality

**Last Updated:** 2025-08-21  
**Version:** 1.1  

## Overview

The preview functionality allows you to quickly view PDF outputs for individual songs in your projects. The system automatically uses your project's ABC file directory setting to find and display available PDFs.

## How to Use Preview

### 1. Individual Song Preview

Each song in your project now has its own **Preview** button:

1. Navigate to a project (Projects → Select Project)
2. In the "Songs in Project" section, find the song you want to preview
   - **Note:** Songs are automatically sorted alphabetically by title for easy navigation
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

## PDF Directory Organization

### Configuring PDF Output Folders

When building projects with Zupfnoter, you can automatically organize different PDF types into separate folders using **Folder Patterns**.

### Setting Up Folder Patterns

1. Go to Projects → Select Project → **Edit Configuration**
2. Add or modify the `folderPatterns` section in the JSON configuration:

```json
{
  "folderPatterns": {
    "*_-*_a3.pdf": "klein",
    "*_-*_a4.pdf": "gross", 
    "*_partitur.pdf": "partituren",
    "*_sammlung*.pdf": "sammlungen",
    "*.pdf": "alle_anderen"
  }
}
```

### How Folder Patterns Work

- **Pattern Matching:** PDF filenames are matched against patterns using wildcards (`*`)
- **First Match Wins:** The first matching pattern determines the target folder
- **Order Matters:** Place more specific patterns before general ones
- **Automatic Creation:** Folders are created automatically during project builds

### Common Pattern Examples

| Pattern | Matches | Goes to Folder |
|---------|---------|----------------|
| `*_klein.pdf` | `song_klein.pdf` | `einzelstimmen_klein/` |
| `*_gross.pdf` | `song_gross.pdf` | `einzelstimmen_gross/` |
| `*_a3.pdf` | `melody_a3.pdf` | `a3_format/` |
| `*_a4.pdf` | `melody_a4.pdf` | `a4_format/` |
| `*_partitur*.pdf` | `song_partitur_full.pdf` | `partituren/` |
| `*_sammlung*.pdf` | `christmas_sammlung.pdf` | `sammlungen/` |
| `*.pdf` | Any remaining PDF | `alle_anderen/` |

### Resulting Directory Structure

After a project build, your output directory will be organized like this:

```
/your_output_directory/
├── klein/                 # Small format PDFs
│   ├── song1_-1_a3.pdf
│   └── song2_-2_a3.pdf
├── gross/                 # Large format PDFs  
│   ├── song1_-1_a4.pdf
│   └── song2_-2_a4.pdf
├── partituren/            # Score PDFs
│   ├── song1_partitur.pdf
│   └── song2_partitur.pdf
├── sammlungen/            # Collection PDFs
│   └── christmas_sammlung.pdf
└── alle_anderen/          # Catch-all folder
    └── special_format.pdf
```

### Best Practices for Folder Patterns

1. **Start Specific:** Put specific patterns first, general patterns last
2. **Use Descriptive Names:** Choose folder names that clearly indicate content
3. **Include Catch-All:** Always end with `"*.pdf": "default_folder"` 
4. **Test Patterns:** Verify patterns match your expected filenames
5. **Document Patterns:** Keep notes about your folder organization system

### Example Configurations

**Simple Organization:**
```json
{
  "folderPatterns": {
    "*_klein.pdf": "klein",
    "*_gross.pdf": "gross",
    "*.pdf": "sonstige"
  }
}
```

**Detailed Organization:**
```json
{
  "folderPatterns": {
    "*_einzelstimme_klein.pdf": "einzelstimmen/klein",
    "*_einzelstimme_gross.pdf": "einzelstimmen/gross",
    "*_partitur_*.pdf": "partituren/vollpartitur", 
    "*_partitur.pdf": "partituren/einfach",
    "*_sammlung_*.pdf": "sammlungen",
    "*_probe*.pdf": "probeblätter",
    "*.pdf": "verschiedenes"
  }
}
```

### Troubleshooting Folder Patterns

**PDFs in Wrong Folders:**
- Check pattern order (specific before general)
- Verify wildcard placement matches actual filenames
- Test with a small build first

**Folders Not Created:**
- Ensure patterns actually match generated filenames
- Check that project build completed successfully
- Verify output directory permissions

**Missing PDFs:**
- Check if patterns are too restrictive
- Ensure catch-all pattern `"*.pdf"` exists
- Review Zupfnoter build logs for errors

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

## Song Organization

### Alphabetical Sorting

**New in Version 1.1:** Songs in the project view are automatically sorted alphabetically by title:

- **Automatic:** No manual sorting required
- **Case-insensitive:** "Apple" and "apple" sort together
- **German locale:** Proper handling of umlauts (ä, ö, ü) and ß
- **Consistent:** Order remains the same regardless of when songs were added
- **Real-time:** New songs automatically appear in correct alphabetical position

### Finding Songs Quickly

With alphabetical sorting, you can:
- **Scan efficiently:** Jump to approximate location based on first letter
- **Predict location:** Know where to look for specific song titles
- **Navigate large projects:** Easily manage projects with many songs
- **Maintain organization:** Professional, organized appearance

### Song Display Order

Songs are sorted by:
1. **Primary:** Song title (alphabetical, A-Z)
2. **Fallback:** Songs without titles appear as "Unknown Song"
3. **Locale:** German character ordering (ä after a, ö after o, etc.)

**Examples of sorting order:**
- "Alle meine Entchen"
- "Amazing Grace" 
- "Ännchen von Tharau"
- "Beethoven Sonata"
- "Über den Wolken"
- "Unknown Song" (for songs without titles)
