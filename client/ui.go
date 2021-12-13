package main

import (
	"image"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype/truetype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/blizzy78/ebitenui"
	ebitenuiimage "github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"

	"golang.org/x/image/font"
)

type LoginUI struct {
	account            string
	password           string
	toLogin            bool
	ui                 *ebitenui.UI
	usernameInputImage *widget.TextInputImage
	passwordInputImage *widget.TextInputImage
	buttonImage        *widget.ButtonImage
	face               font.Face
}

func (l *LoginUI) Init() error {
	var err error

	// load text input image
	l.usernameInputImage, err = loadTextInputImageWithPath("graphics/username.png")
	if err != nil {
		return err
	}

	l.passwordInputImage, err = loadTextInputImageWithPath("graphics/password.png")
	if err != nil {
		return err
	}

	// load button image
	l.buttonImage, err = loadButtonImage()
	if err != nil {
		return err
	}

	// load button text font
	l.face, err = loadFont("fonts/NotoSans-Regular.ttf", 20)
	if err != nil {
		return err
	}

	// 帐号输入框
	accountTextInput := newTextInput(
		l.face,
		widget.TextInputOpts.Image(l.usernameInputImage),
		widget.TextInputOpts.Placeholder("Account"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {}),
	)

	// 密码输入框
	passwordTextInput := newTextInput(
		l.face,
		widget.TextInputOpts.Image(l.passwordInputImage),
		widget.TextInputOpts.Placeholder("Password"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {}),
	)

	// 登录按钮
	enterBtn := newButton(
		// set general widget options
		widget.ButtonOpts.Image(l.buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Login", l.face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			l.account = accountTextInput.InputText
			l.password = passwordTextInput.InputText
			l.toLogin = true
		}),
	)

	/*layout := widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Padding(widget.Insets{
			Top:    250,
			Left:   570,
			Right:  100,
			Bottom: 400,
		}),
		widget.RowLayoutOpts.Spacing(30),
	)*/

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(ebitenuiimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		//widget.ContainerOpts.Layout(layout),
	)
	accountTextInput.SetLocation(image.Rect(500, 300, 760, 336))
	rootContainer.AddChild(accountTextInput)
	passwordTextInput.SetLocation(image.Rect(500, 370, 760, 406))
	rootContainer.AddChild(passwordTextInput)
	enterBtn.SetLocation(image.Rect(585, 440, 675, 481))
	rootContainer.AddChild(enterBtn)

	l.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	return nil
}

func (l *LoginUI) Close() {
	l.face.Close()
}

func (l *LoginUI) Update() {
	l.ui.Update()
}

func (l *LoginUI) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	l.ui.Draw(screen)
}

func (l *LoginUI) ResetLoginState() {
	l.toLogin = false
}

func (l *LoginUI) GetLoginState() (string, string, bool) {
	return l.account, l.password, l.toLogin
}

type UIManager struct {
	owner   IGame
	loginUI *LoginUI
}

func NewUIMgr(owner IGame) *UIManager {
	return &UIManager{
		owner:   owner,
		loginUI: &LoginUI{},
	}
}

func (ui *UIManager) Init() {
	ui.loginUI.Init()
}

func (u *UIManager) Update() {
	u.checkAndLogin()
	u.loginUI.Update()
}

func (u *UIManager) Draw(screen *ebiten.Image) {
	if u.owner.GetState() == GameStateMainMenu {
		u.loginUI.Draw(screen)
	}
}

func (u *UIManager) checkAndLogin() {
	acc, passwd, o := u.loginUI.GetLoginState()
	if o {
		u.owner.EventMgr().InvokeEvent(EventIdOpLogin, acc, passwd)
		u.loginUI.ResetLoginState()
	}
}

func newTextInput(face font.Face, opts ...widget.TextInputOpt) *widget.TextInput {
	ti := widget.NewTextInput(append(opts, []widget.TextInputOpt{
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:     color.Gray16{50},
			Disabled: color.Gray{},
			Caret:    color.Black,
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 1)),
	}...)...)
	//event.ExecuteDeferred()
	//render(ti, ti.GetWidget().Rect.Dx(), ti.GetWidget().Rect.Dy())
	return ti
}

func newButton(opts ...widget.ButtonOpt) *widget.Button {
	b := widget.NewButton(opts...)
	//event.ExecuteDeferred()
	//render(b, b.GetWidget().Rect.Dx(), b.GetWidget().Rect.Dy())
	return b
}

/*
func loadTextInputImage() (*widget.TextInputImage, error) {
	idle, err := loadNineSlice("graphics/text-input-idle.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	disabled, err := loadNineSlice("graphics/text-input-disabled.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	return &widget.TextInputImage{Idle: idle, Disabled: disabled}, nil
}
*/

func loadTextInputImageWithPath(imagePath string) (*widget.TextInputImage, error) {
	img, err := loadNineSlice(imagePath, [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	return &widget.TextInputImage{Idle: img, Disabled: img}, nil
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle, err := loadNineSlice("graphics/button-idle.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	hover, err := loadNineSlice("graphics/button-hover.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	pressed, err := loadNineSlice("graphics/button-pressed.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadNineSlice(path string, w [3]int, h [3]int) (*ebitenuiimage.NineSlice, error) {
	i, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	return ebitenuiimage.NewNineSlice(i, w, h), nil
}

func loadFont(path string, size float64) (font.Face, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
