// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package imageutil

import (
	"bytes"
	"context"
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	"openpitrix.io/openpitrix/pkg/logger"
)

var ErrDecodeImage = fmt.Errorf("failed to decode image")
var ErrEncodeImage = fmt.Errorf("failed to encode image")

func Thumbnail(ctx context.Context, b []byte) (map[string][]byte, error) {
	output := make(map[string][]byte)
	buf := bytes.NewBuffer(b)
	_, ext, err := image.DecodeConfig(buf)
	if err != nil {
		logger.Error(ctx, "Failed to decode image: %+v", err)
		return nil, ErrDecodeImage
	}
	buf = bytes.NewBuffer(b)
	format, _ := imaging.FormatFromExtension(ext)
	img, err := imaging.Decode(buf)
	if err != nil {
		logger.Error(ctx, "Failed to decode image: %+v", err)
		return nil, ErrDecodeImage
	}
	for _, size := range []int{64, 384} {
		dst, err := thumbnail(ctx, format, img, size)
		if err != nil {
			return nil, err
		}
		output[fmt.Sprint(size)] = dst
	}
	output["raw"] = b
	return output, nil
}

func thumbnail(ctx context.Context, format imaging.Format, img image.Image, size int) ([]byte, error) {
	dstimg := imaging.Fit(img, size, size, imaging.Lanczos)
	outBuf := new(bytes.Buffer)
	err := imaging.Encode(outBuf, dstimg, format)
	if err != nil {
		logger.Error(ctx, "Failed to encode image: %+v", err)
		return nil, ErrEncodeImage
	}
	return outBuf.Bytes(), nil
}
