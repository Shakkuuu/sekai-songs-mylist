import { IMAGE_BASE_URL } from "../lib/constants";

const res = await fetch(`${IMAGE_BASE_URL}/upload/thumbnail`, {
  method: "POST",
  body: formData,
});
