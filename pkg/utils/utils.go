package utils

func ChangeNToNormal(kind string) string {
	if kind == "a0_n" {
		return "a0"
	}
	if kind == "b_n" {
		return "b"
	}
	if kind == "f_n" {
		return "f"
	}

	return kind
}
