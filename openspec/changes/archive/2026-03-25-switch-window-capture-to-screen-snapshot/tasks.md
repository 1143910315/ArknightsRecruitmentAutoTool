## 1. Capture Strategy

- [x] 1.1 Replace the current direct window-capture implementation with a full-screen snapshot capture path in the backend screenshot service.
- [x] 1.2 Keep the existing window screenshot entrypoint and return shape stable so recognition-setting flows continue to call the same API.

## 2. Crop Handling

- [x] 2.1 Crop the target window image from the full-screen snapshot using the target window rectangle coordinates.
- [x] 2.2 Clamp the crop rectangle to the captured screen bounds and return a clear failure when the window has no visible intersection.

## 3. Verification

- [x] 3.1 Verify the recognition-setting screenshot flow still shows the captured image after the capture-source switch.
- [x] 3.2 Validate normal and partially out-of-bounds windows to confirm cropped results follow the visible screen area.
