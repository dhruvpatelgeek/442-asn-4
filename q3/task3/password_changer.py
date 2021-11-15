import hashlib

with open('t3.exe', 'rb') as f:
    hexdata = f.read().hex()
flagpole ="396d4d70353266734e24" # find the location where the SHA-1 of the password is stored
loc=hexdata.find(flagpole)+len(flagpole)
print("HASH of the old password is  <"+hexdata[loc:loc+40]+">");


new_password = input('choose a new password: ')
print("new password is : <"+ new_password+">")
hash_object = hashlib.sha1(new_password.encode())
new_password_hash=hash_object.hexdigest()
print("HASH of the NEW password is  <"+new_password_hash+">");

#since strings are immutable in python im just recrating the entire program 
#it is a very inefficient way to do this 
slice_a=hexdata[:loc]
slice_b=new_password_hash
slice_c=hexdata[loc+40:]

new_program=slice_a+slice_b+slice_c
result = bytearray.fromhex(new_program)

print("[password changed]")
flagpole ="396d4d70353266734e24" # find the location where the SHA-1 of the password is stored
loc=new_program.find(flagpole)+len(flagpole)
print("HASH of the CHANGED password <"+new_program[loc:loc+40]+">");
print("[writing output.exe]")

open('Output.exe', 'wb').write(result)
