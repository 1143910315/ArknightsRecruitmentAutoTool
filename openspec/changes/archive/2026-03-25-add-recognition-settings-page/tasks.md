## 1. Recognition Settings Page

- [x] 1.1 Add a recognition settings navigation entry and route it to a dedicated frontend page
- [x] 1.2 Build the recognition settings page layout for window selection, window info, screenshot preview, region list, and save actions
- [x] 1.3 Reuse the existing mouse-pointing window selection flow so the page can capture the current target window

## 2. Window Screenshot and Template Data

- [x] 2.1 Add or expose a backend method that captures the selected target window as an image for preview and region selection
- [x] 2.2 Define a persisted window template structure that stores window identity, screenshot metadata, and multiple region entries
- [x] 2.3 Save region reference images and template metadata to a stable local storage location for future recognition use

## 3. Region Annotation Workflow

- [x] 3.1 Implement screenshot-based drawing or selection of multiple regions on the target window image
- [x] 3.2 Save each region with a label and normalized coordinates relative to the window screenshot
- [x] 3.3 Generate and store a reference image for each selected region from the window screenshot crop
- [x] 3.4 Allow users to review, edit, or remove configured regions before saving the template

## 4. Matching Preparation and Verification

- [x] 4.1 Add a matching-oriented backend or data access path that can load saved region templates for future comparison
- [x] 4.2 Define and return clear per-region match or mismatch results when comparing a current region image to its saved reference image
- [x] 4.3 Verify the recognition settings flow for window selection, screenshot preview, multi-region save, and template reload scenarios
- [x] 4.4 Run `wails build` and confirm the app still packages successfully with the new recognition settings page
