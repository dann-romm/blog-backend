package v1

import "github.com/labstack/echo/v4"

func BindAndValidate(c echo.Context, i any) error {
	err := c.Bind(i)
	if err != nil {
		return err
	}

	err = c.Validate(i)
	if err != nil {
		return err
	}

	return nil
}
