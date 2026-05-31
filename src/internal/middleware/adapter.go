package middleware

import (
	"fmt"
	"reflect"

	"silent/src/internal/context/abstract"
	"silent/src/internal/context/adapter"
	"silent/src/internal/provider"

	"github.com/gofiber/fiber/v3"
)

func processServiceArgument(
	argType reflect.Type,
	fiberCtx fiber.Ctx,
	serviceProvider *provider.ServiceProvider,
) (reflect.Value, error) {
	service, err := serviceProvider.Get(argType)

	if err != nil {
		return reflect.Value{}, fmt.Errorf("service not found: %s", argType)
	}

	return reflect.ValueOf(service), nil
}

func processParamsArgument(
	argType reflect.Type,
	fiberCtx fiber.Ctx,
	_ *provider.ServiceProvider,
) (reflect.Value, error) {
	dtoVal := reflect.New(argType).Interface()

	if err := fiberCtx.Bind().All(dtoVal); err != nil {
		return reflect.Value{}, ParamsValidationError{err: err}
	}

	return reflect.ValueOf(dtoVal).Elem(), nil
}

func processDTOArgument(
	argType reflect.Type,
	fiberCtx fiber.Ctx,
	_ *provider.ServiceProvider,
) (reflect.Value, error) {
	dtoPtr := reflect.New(argType.Elem()).Interface()

	if err := fiberCtx.Bind().Body(dtoPtr); err != nil {
		return reflect.Value{}, BodyValidationError{err: err}
	}

	return reflect.ValueOf(dtoPtr), nil
}

var processArgumentFunctions = map[reflect.Kind]func(reflect.Type, fiber.Ctx, *provider.ServiceProvider) (reflect.Value, error){
	reflect.Interface: processServiceArgument,
	reflect.Struct:    processParamsArgument,
}

var processPtrArgumentFunctions = map[reflect.Kind]func(reflect.Type, fiber.Ctx, *provider.ServiceProvider) (reflect.Value, error){
	reflect.Struct: processDTOArgument,
}

func validateHandler(funcType reflect.Type) {
	if funcType.Kind() != reflect.Func {
		panic("Adapt: handlerFunc must be a function")
	}

	if funcType.NumIn() < 1 {
		panic("Adapt: function must have at least one argument (*HandlerContext)")
	}
	firstArgType := funcType.In(0)
	expectedType := reflect.TypeOf((*abstract.HandlerContext)(nil)).Elem()
	if firstArgType != expectedType {
		panic(fmt.Sprintf("Adapt: first argument must be HandlerContext, got %v", firstArgType))
	}

	if funcType.NumOut() != 2 {
		panic("Adapt: function must return (interface{}, error)")
	}
	if funcType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("Adapt: second return type must be error")
	}
}

func getArguments(funcType reflect.Type, fiberCtx fiber.Ctx, serviceProvider *provider.ServiceProvider) ([]reflect.Value, error) {
	ctx := &adapter.FiberCtxAdapter{Ctx: fiberCtx}
	args := make([]reflect.Value, funcType.NumIn())
	args[0] = reflect.ValueOf(ctx)

	var err error

	for i := 1; i < funcType.NumIn(); i++ {
		argType := funcType.In(i)

		funcMap := processArgumentFunctions
		argKey := argType.Kind()

		if argType.Kind() == reflect.Ptr {
			funcMap = processPtrArgumentFunctions
			argKey = argType.Elem().Kind()
		}

		if processFunc, ok := funcMap[argKey]; ok {
			args[i], err = processFunc(argType, fiberCtx, serviceProvider)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported argument type: %v", argType)
		}
	}

	return args, nil
}

func wrapError(errArg any, fiberCtx fiber.Ctx) error {
	err := errArg.(error)
	if sc, ok := err.(interface{ StatusCode() int }); ok {
		return fiberCtx.Status(sc.StatusCode()).JSON(fiber.Map{"error": err.Error()})
	}
	return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func Adapt(handlerFunc interface{}, serviceProvider *provider.ServiceProvider) fiber.Handler {
	funcValue := reflect.ValueOf(handlerFunc)
	funcType := funcValue.Type()

	validateHandler(funcType)

	return func(fiberCtx fiber.Ctx) error {
		args, err := getArguments(funcType, fiberCtx, serviceProvider)

		if err != nil {
			return wrapError(err, fiberCtx)
		}

		results := funcValue.Call(args)

		data := results[0].Interface()
		errVal := results[1].Interface()

		if errVal != nil {
			return wrapError(err, fiberCtx)
		}

		if data == nil {
			return fiberCtx.SendStatus(fiber.StatusNoContent)
		}

		return fiberCtx.JSON(data)
	}
}
