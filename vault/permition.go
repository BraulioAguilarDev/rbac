package vault

func (w *Wrapper) CanDelete(path string) bool {
	if _, err := w.Client.Logical().Delete(path); err != nil {
		return false
	}

	return true
}

func (w *Wrapper) CanRead(path string) bool {
	if _, err := w.Client.Logical().Read(path); err != nil {
		return false
	}

	return true
}

func (w *Wrapper) CanWrite(path string) bool {
	// The v2 of kv secret engine needs this
	data := make(map[string]interface{})
	info := map[string]string{
		"test": "test",
	}

	data["data"] = info

	if _, err := w.Client.Logical().Write(path, data); err != nil {
		return false
	}

	return true
}
