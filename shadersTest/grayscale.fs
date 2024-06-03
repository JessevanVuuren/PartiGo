// grayscale.fs
#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
uniform sampler2D texture0;
out vec4 finalColor;

void main()
{
    vec4 texelColor = texture(texture0, fragTexCoord);
    float gray = (texelColor.r + texelColor.g + texelColor.b) / 3.0;
    finalColor = vec4(vec3(gray), texelColor.a);
}