package logger

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func DumpWithOptions(showReq, showResp, showBody, showHeaders bool, cb func(string)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var out strings.Builder
		request := c.Request()
		response := c.Response()

		// 🔹 Dump request
		if showReq {
			out.WriteString("--- Request ---\n")
			out.WriteString(fmt.Sprintf("%s %s\n", c.Method(), c.OriginalURL()))

			if showHeaders {
				request.Header.VisitAll(func(k, v []byte) {
					out.WriteString(fmt.Sprintf("%s: %s\n", k, v))
				})
			}

			if showBody {
				body := c.Body()
				if len(body) > 0 {
					out.WriteString("\n" + string(body) + "\n")
				}
			}
		}

		err := c.Next()

		if showResp {
			out.WriteString("\n--- Response ---\n")
			out.WriteString(fmt.Sprintf("Status: %d\n", response.StatusCode()))

			if showHeaders {
				response.Header.VisitAll(func(k, v []byte) {
					out.WriteString(fmt.Sprintf("%s: %s\n", k, v))
				})
			}

			if showBody {
				body := response.Body()
				if len(body) > 0 {
					out.WriteString("\n" + string(body) + "\n")
				}
			}
		}

		if cb != nil {
			cb(out.String())
		} else {
			fmt.Println(out.String())
		}
		return err
	}
}
