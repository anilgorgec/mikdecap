package tzsp

func Parse(pkg []byte, out []byte) (int, error) {
	err := headerValidate(pkg)
	if err != nil {
		return -1, err
	}
	pkg = pkg[hLen:]
	var tagType TagType
	var ln int
	for len(pkg) > 0 {
		tagType, ln, err = parseTag(pkg)
		if err != nil {
			return -1, err
		}
		pkg = pkg[ln:]
		if tagType == TagEnd {
			break
		}
	}
	if tagType != TagEnd {
		return -1, ErrMissingEndTag
	}
	copy(out, pkg)
	return len(pkg), err
}

func headerValidate(pkg []byte) error {
	if len(pkg) < hLen {
		return ErrHeaderTooShort
	}
	if uint8(pkg[0]) != 1 {
		return ErrUnsupportedVersion
	}
	if Type(pkg[1]) != TypeReceivedTagList {
		return ErrUnsupportedPacketType
	}
	return nil
}

func parseTag(pkg []byte) (TagType, int, error) {
	tType := TagType(pkg[0])
	if tType == TagPadding || tType == TagEnd {
		return tType, 1, nil
	}
	if len(pkg) < 2 {
		return 0, 0, ErrTruncatedTag
	}
	ln := int(pkg[1] + 2)

	if len(pkg) < ln {
		return 0, 0, ErrTruncatedTag
	}
	return tType, ln, nil
}
